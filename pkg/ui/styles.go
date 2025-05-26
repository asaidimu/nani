package ui

import "github.com/charmbracelet/lipgloss"

var (
	HistoryStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 1) // Top/Bottom padding 1, Left/Right padding 1

	PromptStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#04B575")).
		Padding(1, 1) // Top/Bottom padding 1, Left/Right padding 1

	PreviewStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF6B6B")).
		Padding(0, 1) // Top/Bottom padding 1, Left/Right padding 1

	TitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1). // Top/Bottom padding 0, Left/Right padding 1
		Bold(true)

	UserMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true)

	AIMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true)
)
