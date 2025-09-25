package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	modeChoices []string
	cursor      int
	selected    bool
	chosenMode  string
}

func New() Model {
	return Model{
		modeChoices: []string{"Encrypt", "Decrypt"},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.modeChoices)-1 {
				m.cursor++
			}
		case "enter":
			m.chosenMode = m.modeChoices[m.cursor]
			m.selected = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := "=== Text Encryption Tool ===\n\n"

	for i, choice := range m.modeChoices {
		if m.cursor == i {
			s += fmt.Sprintf("> %s\n", choice)
		} else {
			s += fmt.Sprintf("  %s\n", choice)
		}
	}

	s += "\nUse ↑/↓ arrows to navigate, Enter to select, q to quit"
	return s
}

func (m Model) GetSelectedMode() (string, bool) {
	return m.chosenMode, m.selected
}
