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
	loading      bool
	err          error
	showOnlyFail bool
}

// initialModel creates a new model with default values
func initialModel(testArgs []string) model {
	return model{
		tests:        []TestResult{},
		cursor:       0,
		sidebarWidth: 30,
		ready:        false,
		testArgs:     testArgs,
		loading:      true,
		err:          nil,
		showOnlyFail: false,
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
		m.loading = false
		m.err = msg.err
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

		case "f":
			m.showOnlyFail = !m.showOnlyFail
			m.cursor = 0 // Reset cursor when toggling filter
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
		return "Initializing..."
	}

	if m.loading {
		return loadingStyle.Render("Running tests...") + "\n\n" +
			dimStyle.Render("This may take a moment depending on your test suite")
	}

	if m.err != nil {
		return errorStyle.Render("Error: "+m.err.Error()) + "\n\n" +
			dimStyle.Render("Press q to quit")
	}

	if len(m.tests) == 0 {
		return "No tests found.\n\nPress q to quit."
	}

	sidebar := m.renderSidebar()
	content := m.renderContent()
	statusBar := m.renderStatusBar()

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
	return lipgloss.JoinVertical(lipgloss.Left, mainView, statusBar)
}

// renderSidebar builds the sidebar with the test list
func (m model) renderSidebar() string {
	var sidebarContent strings.Builder

	// Title with filter indicator
	title := "Tests"
	if m.showOnlyFail {
		title += " (Failures Only)"
	}
	sidebarContent.WriteString(titleStyle.Render(title) + "\n\n")

	// Get filtered tests
	filteredTests := m.getFilteredTests()

	if len(filteredTests) == 0 {
		if m.showOnlyFail {
			sidebarContent.WriteString(successStyle.Render("✓ All tests passed!") + "\n")
		} else {
			sidebarContent.WriteString("No tests to display\n")
		}
	} else {
		visibleStart, visibleEnd := m.calculateVisibleRange()

		for i := visibleStart; i < visibleEnd && i < len(filteredTests); i++ {
			test := filteredTests[i]
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
	}

	return sidebarStyle.
		Width(m.sidebarWidth - 2).
		Height(m.height - 4).
		Render(sidebarContent.String())
}

// renderContent builds the content pane showing the selected test details
func (m model) renderContent() string {
	var contentText strings.Builder

	filteredTests := m.getFilteredTests()

	if len(filteredTests) == 0 {
		// Show success message when all tests pass
		if m.showOnlyFail && m.hasPassingTests() {
			contentText.WriteString(successStyle.Render("🎉 All Tests Passed!") + "\n\n")
			contentText.WriteString("All tests in your test suite passed successfully.\n")
			contentText.WriteString(dimStyle.Render("Press 'f' to show all tests"))
		}
	} else if m.cursor < len(filteredTests) {
		selectedTest := filteredTests[m.cursor]

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
		Height(m.height - 4).
		Render(contentText.String())
}

// calculateVisibleRange determines which tests should be visible in the sidebar
func (m model) calculateVisibleRange() (int, int) {
	filteredTests := m.getFilteredTests()
	visibleStart := 0
	visibleEnd := len(filteredTests)

	// Calculate visible range to fit in window
	maxVisible := m.height - 6 // Account for title, borders, padding, and status bar
	if len(filteredTests) > maxVisible {
		visibleStart = max(m.cursor-maxVisible/2, 0)
		visibleEnd = visibleStart + maxVisible
		if visibleEnd > len(filteredTests) {
			visibleEnd = len(filteredTests)
			visibleStart = max(visibleEnd-maxVisible, 0)
		}
	}

	return visibleStart, visibleEnd
}

// getFilteredTests returns tests based on the current filter
func (m model) getFilteredTests() []TestResult {
	if !m.showOnlyFail {
		return m.tests
	}

	filtered := []TestResult{}
	for _, test := range m.tests {
		if test.Status == "fail" {
			filtered = append(filtered, test)
		}
	}
	return filtered
}

// hasPassingTests checks if there are any passing tests
func (m model) hasPassingTests() bool {
	for _, test := range m.tests {
		if test.Status == "pass" {
			return true
		}
	}
	return false
}

// renderStatusBar renders the status bar at the bottom
func (m model) renderStatusBar() string {
	failCount := 0
	passCount := 0
	for _, test := range m.tests {
		if test.Status == "fail" {
			failCount++
		} else if test.Status == "pass" {
			passCount++
		}
	}

	var statusParts []string

	if passCount > 0 {
		statusParts = append(statusParts, passStyle.Render(fmt.Sprintf("%d passed", passCount)))
	}
	if failCount > 0 {
		statusParts = append(statusParts, failStyle.Render(fmt.Sprintf("%d failed", failCount)))
	}

	status := strings.Join(statusParts, " • ")

	filterHint := dimStyle.Render("f: toggle filter")
	quitHint := dimStyle.Render("q: quit")

	leftSide := status
	rightSide := filterHint + " • " + quitHint

	gap := m.width - lipgloss.Width(leftSide) - lipgloss.Width(rightSide)
	if gap < 0 {
		gap = 0
	}

	statusLine := leftSide + strings.Repeat(" ", gap) + rightSide

	return statusBarStyle.
		Width(m.width).
		Render(statusLine)
}
