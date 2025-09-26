package main

import (
	"fmt"
	"os"
	"txt-encdec-cli/tui"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	appName    = "Text Encryptor"
	appVersion = "2.0.0"
)

func main() {
	model := tui.New()
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s v%s: %v\n", appName, appVersion, err)
		os.Exit(1)
	}
}
