package tui

import (
	"fmt"
	"txt-encdec-cli/core"
	"txt-encdec-cli/platform"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state        AppState
	mode         OperationMode
	cursor       int
	terminalSize TerminalSize

	textInput  textinput.Model
	inputState InputState

	cryptor   core.Cryptor
	clipboard platform.ClipboardManager
	detector  platform.SystemStateDetector

	layout *LayoutManager
	config AppConfig

	secretKey string
	result    string
	lastError error

	availableModes []string
}

func New() Model {
	return NewWithConfig(DefaultConfig())
}

func NewWithConfig(config AppConfig) Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = config.InputCharLimit

	return Model{
		state:          StateSelectMode,
		terminalSize:   TerminalSize{Width: config.DefaultWidth, Height: config.DefaultHeight},
		textInput:      ti,
		clipboard:      platform.NewLinuxClipboardManager(),
		detector:       platform.NewLinuxSystemDetector(),
		layout:         NewLayoutManager(config),
		config:         config,
		availableModes: []string{"Encrypt", "Decrypt"},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalSize = TerminalSize{Width: msg.Width, Height: msg.Height}

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		if m.state == StateEnterSecret {
			m.updateInputState(msg)
		} else {
			m.clearInputState()
		}

		if cmd := m.handleKeyEvent(msg); cmd != nil {
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m *Model) updateInputState(msg tea.KeyMsg) {
	m.inputState.CapsLockOn = m.detector.IsCapsLockOn()

	if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
		hasKorean := false
		hasLatin := false

		for _, r := range msg.Runes {
			if m.detector.IsKoreanInput(r) {
				hasKorean = true
			} else if m.detector.IsLatinInput(r) {
				hasLatin = true
			}
		}

		if hasKorean {
			m.inputState.KoreanActive = true
		} else if hasLatin {
			m.inputState.KoreanActive = false
		}
	}
}

func (m *Model) clearInputState() {
	m.inputState = InputState{}
}

func (m *Model) handleKeyEvent(msg tea.KeyMsg) tea.Cmd {
	switch m.state {
	case StateSelectMode:
		return m.handleModeSelection(msg)
	case StateEnterSecret:
		return m.handleSecretEntry(msg)
	case StateEnterText:
		return m.handleTextEntry(msg)
	case StateShowResult, StateShowError:
		return m.handleResultScreen(msg)
	}
	return nil
}

func (m *Model) handleModeSelection(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q":
		return tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.availableModes)-1 {
			m.cursor++
		}
	case "enter":
		m.mode = OperationMode(m.cursor)
		m.transitionToSecretEntry()
		return textinput.Blink
	}
	return nil
}

func (m *Model) handleSecretEntry(msg tea.KeyMsg) tea.Cmd {
	if msg.Type == tea.KeyEnter {
		m.secretKey = m.textInput.Value()
		m.cryptor = core.NewAESCryptor(m.secretKey)
		m.transitionToTextEntry()
		return textinput.Blink
	}
	return nil
}

func (m *Model) handleTextEntry(msg tea.KeyMsg) tea.Cmd {
	if msg.Type == tea.KeyEnter {
		inputText := m.textInput.Value()
		m.processInput(inputText)
	}
	return nil
}

func (m *Model) handleResultScreen(msg tea.KeyMsg) tea.Cmd {
	if msg.Type == tea.KeyEnter {
		return m.resetToModeSelection()
	}
	return nil
}

func (m *Model) transitionToSecretEntry() {
	m.state = StateEnterSecret
	m.textInput.Prompt = ""
	m.textInput.EchoMode = textinput.EchoPassword
	m.textInput.EchoCharacter = '*'
	m.textInput.Reset()
}

func (m *Model) transitionToTextEntry() {
	m.state = StateEnterText
	m.textInput.Prompt = ""
	m.textInput.EchoMode = textinput.EchoNormal
	m.textInput.Reset()
}

func (m *Model) processInput(inputText string) {
	var result string
	var err error

	switch m.mode {
	case ModeEncrypt:
		result, err = m.cryptor.Encrypt(inputText)
	case ModeDecrypt:
		result, err = m.cryptor.Decrypt(inputText)
	default:
		err = &AppError{Op: "process_input", Err: ErrInvalidOperation}
	}

	if err != nil {
		m.state = StateShowError
		m.lastError = err
	} else {
		m.state = StateShowResult
		m.result = result
		_ = m.clipboard.Copy(result)
	}
}

func (m *Model) resetToModeSelection() tea.Cmd {
	newModel := NewWithConfig(m.config)
	newModel.terminalSize = m.terminalSize
	*m = newModel
	return nil
}

func (m Model) View() string {
	var content string

	switch m.state {
	case StateSelectMode:
		content = m.layout.RenderModeSelection(m.cursor, m.availableModes)

	case StateEnterSecret:
		inputWidth := m.layout.CalculateInputWidth(m.terminalSize)
		m.textInput.Width = inputWidth
		inputView := m.layout.CreateStyledInput(m.textInput.View(), inputWidth)
		content = m.layout.RenderInputPrompt("Enter Secret Key:", inputView, "enter: confirm , ctrl+c: quit")
		content += m.layout.RenderInputState(m.inputState)

	case StateEnterText:
		inputWidth := m.layout.CalculateInputWidth(m.terminalSize)
		m.textInput.Width = inputWidth
		inputView := m.layout.CreateStyledInput(m.textInput.View(), inputWidth)
		title := fmt.Sprintf("Enter Text to %s:", m.mode.String())
		content = m.layout.RenderInputPrompt(title, inputView, "enter: confirm , ctrl+c: quit")

	case StateShowResult:
		message := "Success! Result copied to clipboard"
		details := fmt.Sprintf("Result length: %d characters", len(m.result))
		content = m.layout.RenderResult(true, message, details)

	case StateShowError:
		message := fmt.Sprintf("Error: %v", m.lastError)
		content = m.layout.RenderResult(false, message, "")
	}

	return m.layout.RenderApp(content)
}
