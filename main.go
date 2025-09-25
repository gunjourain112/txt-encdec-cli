package main

import (
	"fmt"
	"os"
	"txt-encdec-cli/platform"
	"txt-encdec-cli/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	defer platform.RestoreTerminal()

	m := tui.New()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
