package main

import (
	"fmt"
	"os"
	"txt-encdec-cli/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := tui.New()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
