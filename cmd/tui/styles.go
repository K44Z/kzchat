package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Underline(true)
	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	InputStyle = lipgloss.NewStyle().
			Padding(0, 1).
			BorderStyle(lipgloss.HiddenBorder())

	FocusedStyle = InputStyle.Copy().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("205"))
)

var layoutStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Margin(0)

func StatusStyle(mode string) lipgloss.Style {
	var background string
	switch mode {
	case "Insert":
		background = "#85c1dc"
	case "Command":
		background = "#babbf1"
	case "Normal":
		background = "#585b70"
	case "Search":
		background = "#c6a0f6"
	default:
		background = "#b7bdf8"
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#eff1f5")).
		Background(lipgloss.Color(background)).
		Padding(0, 1).
		Bold(true).BorderRight(true).
		BorderForeground(lipgloss.Color("#89b4fa"))
}

var statusRight = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#eff1f5")).
	Background(lipgloss.Color("#8839ef")).
	Padding(0, 1)

var statusBar = lipgloss.NewStyle().
	Background(lipgloss.Color("#414559")).
	Height(1).
	Width(0)

var statusMid = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	Background(lipgloss.Color("#333333")).Align(lipgloss.Center)

var commandStyle = lipgloss.NewStyle()
