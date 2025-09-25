package tui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor = lipgloss.Color("#7C3AED")
	successColor = lipgloss.Color("#10B981")
	errorColor   = lipgloss.Color("#EF4444")
	mutedColor   = lipgloss.Color("#6B7280")
	accentColor  = lipgloss.Color("#F59E0B")

	AppStyle = lipgloss.NewStyle().
			Padding(2, 4).
			Margin(1, 2)

	LogoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1)

	ListPromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			MarginBottom(1)

	ListItemStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginBottom(0)

	SelectedListItemStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				MarginBottom(0)

	TextInputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginBottom(1)

	ResultStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			MarginBottom(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			MarginBottom(1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	CapsIndicatorStyle = lipgloss.NewStyle().
				Background(errorColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 1).
				Bold(true).
				MarginRight(1)

	KoreanIndicatorStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#3B82F6")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 1).
				Bold(true).
				MarginRight(1)
)
