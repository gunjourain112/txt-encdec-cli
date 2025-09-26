package tui

import (
	"strings"
)

type LayoutManager struct {
	config AppConfig
}

func NewLayoutManager(config AppConfig) *LayoutManager {
	return &LayoutManager{
		config: config,
	}
}

func (lm *LayoutManager) CalculateInputWidth(terminalSize TerminalSize) int {
	if !terminalSize.IsValid() || terminalSize.Width < lm.config.MinTerminalWidth {
		return lm.config.MinInputWidth
	}

	calculated := terminalSize.Width - 16 

	if calculated > lm.config.MaxInputWidth {
		return lm.config.MaxInputWidth
	}

	if calculated < lm.config.MinInputWidth {
		return lm.config.MinInputWidth
	}

	return calculated
}

func (lm *LayoutManager) RenderModeSelection(cursor int, modes []string) string {
	var content strings.Builder

	content.WriteString(ListPromptStyle.Render("Select encryption mode:") + "\n")

	for i, mode := range modes {
		if cursor == i {
			content.WriteString(SelectedListItemStyle.Render("> "+mode) + "\n")
		} else {
			content.WriteString(ListItemStyle.Render("  "+mode) + "\n")
		}
	}

	content.WriteString("\n" + HelpStyle.Render("up/down: navigate , enter: select , q/ctrl+c: quit"))

	return content.String()
}

func (lm *LayoutManager) RenderInputPrompt(title, inputView, helpText string) string {
	var content strings.Builder

	content.WriteString(ListPromptStyle.Render(title) + "\n")
	content.WriteString(inputView + "\n")
	content.WriteString(HelpStyle.Render(helpText))

	return content.String()
}

func (lm *LayoutManager) RenderResult(success bool, message, details string) string {
	var content strings.Builder

	if success {
		content.WriteString(ResultStyle.Render(" "+message) + "\n\n")
	} else {
		content.WriteString(ErrorStyle.Render(" "+message) + "\n\n")
	}

	if details != "" {
		content.WriteString(HelpStyle.Render(details) + "\n\n")
	}

	content.WriteString(HelpStyle.Render("enter: continue"))

	return content.String()
}

func (lm *LayoutManager) RenderInputState(state InputState) string {
	if !state.HasIndicators() {
		return ""
	}

	var indicators []string

	if state.KoreanActive {
		indicators = append(indicators, KoreanIndicatorStyle.Render("한글"))
	}

	if state.CapsLockOn {
		indicators = append(indicators, CapsIndicatorStyle.Render("CAPS"))
	}

	return "\n\n" + strings.Join(indicators, " ")
}

func (lm *LayoutManager) RenderApp(content string) string {
	var app strings.Builder

	logo := " TEXT ENCRYPTOR "
	app.WriteString(LogoStyle.Render(logo) + "\n\n")

	app.WriteString(content)

	return AppStyle.Render(app.String())
}

func (lm *LayoutManager) CreateStyledInput(inputView string, width int) string {
	styleWidth := width + 6 
	inputStyle := TextInputStyle.Width(styleWidth)
	return inputStyle.Render(inputView)
}
