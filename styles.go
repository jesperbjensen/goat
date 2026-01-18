package main

import "github.com/charmbracelet/lipgloss"

// UI style definitions
var (
	// sidebarStyle defines the styling for the test list sidebar
	sidebarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true).
			BorderForeground(lipgloss.Color("240")).
			PaddingRight(1)

	// contentStyle defines the styling for the main content area
	contentStyle = lipgloss.NewStyle().
			Padding(0, 2)

	// selectedItemStyle is used for the currently selected test in the sidebar
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true)

	// passStyle is used for passing tests (green)
	passStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	// failStyle is used for failing tests (red)
	failStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	// titleStyle is used for section titles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))
)
