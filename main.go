package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TestEvent struct {
	Time    string  `json:"Time"`
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

type TestResult struct {
	Name     string
	Status   string
	Output   []string
	FilePath string
	Line     int
}

type model struct {
	tests        []TestResult
	cursor       int
	width        int
	height       int
	sidebarWidth int
	ready        bool
}

type testsLoadedMsg struct {
	tests []TestResult
}

var (
	sidebarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true).
			BorderForeground(lipgloss.Color("240")).
			PaddingRight(1)

	contentStyle = lipgloss.NewStyle().
			Padding(0, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true)

	passStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	failStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))
)

func initialModel() model {
	return model{
		tests:        []TestResult{},
		cursor:       0,
		sidebarWidth: 30,
		ready:        false,
	}
}

func (m model) Init() tea.Cmd {
	return loadTests
}

func loadTests() tea.Msg {
	cmd := exec.Command("go", "test", "-json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return testsLoadedMsg{tests: []TestResult{}}
	}

	if err := cmd.Start(); err != nil {
		return testsLoadedMsg{tests: []TestResult{}}
	}

	scanner := bufio.NewScanner(stdout)
	testResults := make(map[string]*TestResult)

	for scanner.Scan() {
		line := scanner.Bytes()
		var event TestEvent

		if err := json.Unmarshal(line, &event); err != nil {
			continue
		}

		if event.Test != "" && event.Action == "output" {
			if testResults[event.Test] == nil {
				testResults[event.Test] = &TestResult{
					Name:   event.Test,
					Output: []string{},
				}
			}
			testResults[event.Test].Output = append(testResults[event.Test].Output, event.Output)

			if testResults[event.Test].FilePath == "" {
				trimmed := strings.TrimSpace(event.Output)
				if strings.Contains(trimmed, ".go:") {
					parts := strings.Split(trimmed, ":")
					if len(parts) >= 2 {
						testResults[event.Test].FilePath = parts[0]
						if len(parts) >= 2 {
							var lineNum int
							fmt.Sscanf(parts[1], "%d", &lineNum)
							testResults[event.Test].Line = lineNum
						}
					}
				}
			}
		}

		if event.Test != "" && (event.Action == "pass" || event.Action == "fail") {
			if testResults[event.Test] == nil {
				testResults[event.Test] = &TestResult{
					Name:   event.Test,
					Output: []string{},
				}
			}
			testResults[event.Test].Status = event.Action
		}
	}

	cmd.Wait()

	// Convert map to sorted slice
	tests := make([]TestResult, 0, len(testResults))
	for _, result := range testResults {
		tests = append(tests, *result)
	}

	// Sort: failed tests first, then by name
	sort.Slice(tests, func(i, j int) bool {
		if tests[i].Status == "fail" && tests[j].Status != "fail" {
			return true
		}
		if tests[i].Status != "fail" && tests[j].Status == "fail" {
			return false
		}
		return tests[i].Name < tests[j].Name
	})

	return testsLoadedMsg{tests: tests}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case testsLoadedMsg:
		m.tests = msg.tests
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.tests)-1 {
				m.cursor++
			}

		case "home", "g":
			m.cursor = 0

		case "end", "G":
			m.cursor = len(m.tests) - 1

		case "enter":
			if m.cursor < len(m.tests) {
				test := m.tests[m.cursor]
				if test.FilePath != "" && test.Line > 0 {
					return m, openFile(test.FilePath, test.Line)
				}
			}
		}

	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.MouseButtonWheelDown:
			if m.cursor < len(m.tests)-1 {
				m.cursor++
			}
		case tea.MouseButtonLeft:
			// Check if click is in sidebar
			if msg.X < m.sidebarWidth {
				// Calculate which test was clicked (accounting for header)
				clickedIndex := msg.Y - 2
				if clickedIndex >= 0 && clickedIndex < len(m.tests) {
					m.cursor = clickedIndex
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Loading tests..."
	}

	if len(m.tests) == 0 {
		return "No tests found.\n\nPress q to quit."
	}

	// Build sidebar
	var sidebarContent strings.Builder
	sidebarContent.WriteString(titleStyle.Render("Tests") + "\n\n")

	visibleStart := 0
	visibleEnd := len(m.tests)

	// Calculate visible range to fit in window
	maxVisible := m.height - 4 // Account for title, borders, and padding
	if len(m.tests) > maxVisible {
		visibleStart = m.cursor - maxVisible/2
		if visibleStart < 0 {
			visibleStart = 0
		}
		visibleEnd = visibleStart + maxVisible
		if visibleEnd > len(m.tests) {
			visibleEnd = len(m.tests)
			visibleStart = visibleEnd - maxVisible
			if visibleStart < 0 {
				visibleStart = 0
			}
		}
	}

	for i := visibleStart; i < visibleEnd; i++ {
		test := m.tests[i]
		cursor := " "
		if i == m.cursor {
			cursor = "►"
		}

		status := "?"
		style := lipgloss.NewStyle()
		if test.Status == "pass" {
			status = "✓"
			style = passStyle
		} else if test.Status == "fail" {
			status = "✗"
			style = failStyle
		}

		testName := test.Name
		maxNameLen := m.sidebarWidth - 6
		if len(testName) > maxNameLen {
			testName = testName[:maxNameLen-3] + "..."
		}

		line := fmt.Sprintf("%s %s %s", cursor, status, testName)
		if i == m.cursor {
			line = selectedItemStyle.Render(line)
		} else {
			line = style.Render(line)
		}
		sidebarContent.WriteString(line + "\n")
	}

	sidebar := sidebarStyle.
		Width(m.sidebarWidth - 2).
		Height(m.height - 2).
		Render(sidebarContent.String())

	// Build content pane
	var contentText strings.Builder
	if m.cursor < len(m.tests) {
		selectedTest := m.tests[m.cursor]

		// Title
		statusText := ""
		if selectedTest.Status == "pass" {
			statusText = passStyle.Render("PASS")
		} else if selectedTest.Status == "fail" {
			statusText = failStyle.Render("FAIL")
		}

		contentText.WriteString(titleStyle.Render(selectedTest.Name) + " " + statusText + "\n")

		// File link
		if selectedTest.FilePath != "" && selectedTest.Line > 0 {
			fileLink := fmt.Sprintf("%s:%d", selectedTest.FilePath, selectedTest.Line)
			contentText.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render(fileLink) + "\n")
		}
		contentText.WriteString("\n")

		// Output
		for _, line := range selectedTest.Output {
			trimmed := strings.TrimRight(line, "\n")
			if trimmed != "" {
				contentText.WriteString(trimmed + "\n")
			}
		}
	}

	content := contentStyle.
		Width(m.width - m.sidebarWidth - 4).
		Height(m.height - 2).
		Render(contentText.String())

	// Combine sidebar and content
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
}

func openFile(filepath string, line int) tea.Cmd {
	return func() tea.Msg {
		editors := []string{
			"zed",
			"code",
			"subl",
			"atom",
			"vim",
			"nvim",
			"emacs",
		}

		// Try each editor
		for _, editor := range editors {
			var cmd *exec.Cmd
			switch editor {
			case "zed", "code", "subl", "atom":
				cmd = exec.Command(editor, fmt.Sprintf("%s:%d", filepath, line))
			case "vim", "nvim":
				cmd = exec.Command(editor, "+"+fmt.Sprint(line), filepath)
			case "emacs":
				cmd = exec.Command(editor, fmt.Sprintf("+%d", line), filepath)
			}

			if err := cmd.Start(); err == nil {
				return nil
			}
		}

		// Fallback: try to open with default editor
		cmd := exec.Command("open", fmt.Sprintf("%s:%d", filepath, line))
		cmd.Start()

		return nil
	}
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
