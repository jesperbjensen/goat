package main

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// openFile attempts to open a file at a specific line in an available editor
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
