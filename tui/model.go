package tui

import (
	"fmt"
	"strings"
	"txt-encdec-cli/core"
	"txt-encdec-cli/platform"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateChooseMode state = iota
	stateEnterSecret
	stateEnterText
	stateShowResult
	stateShowError
)

type Model struct {
	state       state
	modeChoices []string
	cursor      int
	chosenMode  string
	textInput   textinput.Model
	secretKey   string
	result      string
	err         error
	isCaps      bool
	hasKorean   bool
}

func New() Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 1024
	ti.Width = 46

	return Model{
		state:       stateChooseMode,
		modeChoices: []string{"Encrypt", "Decrypt"},
		textInput:   ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		if m.state == stateEnterSecret {
			m.isCaps = platform.IsCapsLockOn()

			if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
				isKoreanInput := false
				for _, r := range msg.Runes {
					if (r >= 0xAC00 && r <= 0xD7A3) ||
						(r >= 0x3131 && r <= 0x318E) {
						isKoreanInput = true
						break
					}
				}

				if isKoreanInput {
					m.hasKorean = true
				} else {
					m.hasKorean = false
				}
			}
		} else {
			m.isCaps = false
			m.hasKorean = false
		}

		switch m.state {
		case stateChooseMode:
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
				m.state = stateEnterSecret
				m.textInput.Prompt = ""
				m.textInput.EchoMode = textinput.EchoPassword
				m.textInput.EchoCharacter = '*'
				return m, textinput.Blink
			}
		case stateEnterSecret:
			if msg.Type == tea.KeyEnter {
				m.secretKey = m.textInput.Value()
				m.state = stateEnterText
				m.textInput.Prompt = ""
				m.textInput.EchoMode = textinput.EchoNormal
				m.textInput.Reset()
				return m, textinput.Blink
			}
		case stateEnterText:
			if msg.Type == tea.KeyEnter {
				inputText := m.textInput.Value()
				var err error
				if m.chosenMode == "Encrypt" {
					m.result, err = core.Encrypt(m.secretKey, inputText)
				} else {
					m.result, err = core.Decrypt(m.secretKey, inputText)
				}
				if err != nil {
					m.state = stateShowError
					m.err = err
				} else {
					m.state = stateShowResult
					_ = platform.CopyToClipboard(m.result)
				}
				return m, nil
			}
		case stateShowResult, stateShowError:
			if msg.Type == tea.KeyEnter {
				return New(), nil
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var content strings.Builder

	logo := " TEXT ENCRYPTOR "
	content.WriteString(LogoStyle.Render(logo) + "\n\n")

	switch m.state {
	case stateChooseMode:
		content.WriteString(ListPromptStyle.Render("Select encryption mode:") + "\n")
		for i, choice := range m.modeChoices {
			if m.cursor == i {
				content.WriteString(SelectedListItemStyle.Render("> "+choice) + "\n")
			} else {
				content.WriteString(ListItemStyle.Render("  "+choice) + "\n")
			}
		}
		content.WriteString("\n" + HelpStyle.Render("up/down: navigate , enter: select , q/ctrl+c: quit"))

	case stateEnterSecret:
		content.WriteString(ListPromptStyle.Render("Enter Secret Key:") + "\n")
		content.WriteString(TextInputStyle.Render(m.textInput.View()) + "\n")
		content.WriteString(HelpStyle.Render("enter: confirm , ctrl+c: quit"))

	case stateEnterText:
		content.WriteString(ListPromptStyle.Render(fmt.Sprintf("Enter Text to %s:", m.chosenMode)) + "\n")
		content.WriteString(TextInputStyle.Render(m.textInput.View()) + "\n")
		content.WriteString(HelpStyle.Render("enter: confirm , ctrl+c: quit"))

	case stateShowResult:
		content.WriteString(ResultStyle.Render(" Success! Result copied to clipboard") + "\n\n")
		content.WriteString(HelpStyle.Render(fmt.Sprintf("Result length: %d characters", len(m.result))) + "\n\n")
		content.WriteString(HelpStyle.Render("enter: continue"))

	case stateShowError:
		content.WriteString(ErrorStyle.Render(" Error: "+m.err.Error()) + "\n\n")
		content.WriteString(HelpStyle.Render("enter: continue"))
	}

	if m.state == stateEnterSecret {
		var statusParts []string
		if m.hasKorean {
			statusParts = append(statusParts, KoreanIndicatorStyle.Render("한글"))
		}
		if m.isCaps {
			statusParts = append(statusParts, CapsIndicatorStyle.Render("CAPS"))
		}

		if len(statusParts) > 0 {
			content.WriteString("\n\n" + strings.Join(statusParts, " "))
		}
	}

	return AppStyle.Render(content.String())
}
