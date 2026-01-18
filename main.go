package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command-line arguments
	// os.Args[0] is the program name, rest are test arguments
	testArgs := os.Args[1:]

	// If no arguments provided, default to ./...
	if len(testArgs) == 0 {
		testArgs = []string{"./..."}
	}

	p := tea.NewProgram(
		initialModel(testArgs),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
