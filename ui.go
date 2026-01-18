package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// model represents the application state for the TUI
type model struct {
	tests        []TestResult
	cursor       int
	width        int
	height       int
	sidebarWidth int
	ready        bool
	testArgs     []string
}

// initialModel creates a new model with default values
func initialModel(testArgs []string) model {
	return model{
		tests:        []TestResult{},
		cursor:       0,
		sidebarWidth: 30,
		ready:        false,
		testArgs:     testArgs,
	}
}

// Init initializes the model and returns the initial command
func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return loadTests(m.testArgs)
	}
}

// Update handles messages and updates the model
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

// View renders the UI
func (m model) View() string {
	if !m.ready {
		return "Loading tests..."
	}

	if len(m.tests) == 0 {
		return "No tests found.\n\nPress q to quit."
	}

	sidebar := m.renderSidebar()
	content := m.renderContent()

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
}

// renderSidebar builds the sidebar with the test list
func (m model) renderSidebar() string {
	var sidebarContent strings.Builder
	sidebarContent.WriteString(titleStyle.Render("Tests") + "\n\n")

	visibleStart, visibleEnd := m.calculateVisibleRange()

	for i := visibleStart; i < visibleEnd; i++ {
		test := m.tests[i]
		cursor := " "
		if i == m.cursor {
			cursor = "►"
		}

		status := "?"
		style := lipgloss.NewStyle()
		switch test.Status {
		case "pass":
			status = "✓"
			style = passStyle
		case "fail":
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

	return sidebarStyle.
		Width(m.sidebarWidth - 2).
		Height(m.height - 2).
		Render(sidebarContent.String())
}

// renderContent builds the content pane showing the selected test details
func (m model) renderContent() string {
	var contentText strings.Builder

	if m.cursor < len(m.tests) {
		selectedTest := m.tests[m.cursor]

		// Title
		statusText := ""
		switch selectedTest.Status {
		case "pass":
			statusText = passStyle.Render("PASS")
		case "fail":
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

	return contentStyle.
		Width(m.width - m.sidebarWidth - 4).
		Height(m.height - 2).
		Render(contentText.String())
}

// calculateVisibleRange determines which tests should be visible in the sidebar
func (m model) calculateVisibleRange() (int, int) {
	visibleStart := 0
	visibleEnd := len(m.tests)

	// Calculate visible range to fit in window
	maxVisible := m.height - 4 // Account for title, borders, and padding
	if len(m.tests) > maxVisible {
		visibleStart = max(m.cursor-maxVisible/2, 0)
		visibleEnd = visibleStart + maxVisible
		if visibleEnd > len(m.tests) {
			visibleEnd = len(m.tests)
			visibleStart = max(visibleEnd-maxVisible, 0)
		}
	}

	return visibleStart, visibleEnd
}
