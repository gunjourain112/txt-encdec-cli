package tui

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryColor    = lipgloss.Color("#7C3AED")
	SuccessColor    = lipgloss.Color("#10B981")
	ErrorColor      = lipgloss.Color("#EF4444")
	WarningColor    = lipgloss.Color("#F59E0B")
	InfoColor       = lipgloss.Color("#3B82F6")
	MutedColor      = lipgloss.Color("#6B7280")
	BackgroundColor = lipgloss.Color("#1F2937")
	ForegroundColor = lipgloss.Color("#F3F4F6")
	WhiteColor      = lipgloss.Color("#FFFFFF")
	BlackColor      = lipgloss.Color("#000000")
)

const (
	AppPaddingHorizontal = 4
	AppPaddingVertical   = 2
	AppMarginHorizontal  = 2
	AppMarginVertical    = 1
	InputPadding         = 1
	IndicatorPadding     = 1
)

var (
	AppStyle = lipgloss.NewStyle().
			Padding(AppPaddingVertical, AppPaddingHorizontal).
			Margin(AppMarginVertical, AppMarginHorizontal)

	LogoStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)
)

var (
	ListPromptStyle = lipgloss.NewStyle().
			Foreground(WhiteColor).
			Bold(true).
			MarginBottom(1)

	ListItemStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginBottom(0)

	SelectedListItemStyle = lipgloss.NewStyle().
				Foreground(WarningColor).
				Bold(true).
				MarginBottom(0)

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true)
)

var (
	TextInputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(InputPadding, InputPadding*2).
		MarginBottom(1)
)

var (
	ResultStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			MarginBottom(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true).
			MarginBottom(1)

	CodeStyle = lipgloss.NewStyle().
			Background(BackgroundColor).
			Foreground(ForegroundColor).
			Padding(InputPadding, InputPadding*2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(MutedColor)
)

var (
	baseIndicatorStyle = lipgloss.NewStyle().
				Padding(0, IndicatorPadding).
				Bold(true).
				MarginRight(1)

	CapsIndicatorStyle = baseIndicatorStyle.Copy().
				Background(ErrorColor).
				Foreground(WhiteColor)

	KoreanIndicatorStyle = baseIndicatorStyle.Copy().
				Background(InfoColor).
				Foreground(WhiteColor)

	StatusIndicatorStyle = baseIndicatorStyle.Copy().
				Background(WarningColor).
				Foreground(BlackColor)
)
