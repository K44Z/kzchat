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
		background = "4" 
	case "Command":
		background = "3" 
	case "Normal":
		background = "8"
	case "Search":
		background = "5" 
	default:
		background = "8"
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color(background)).
		Padding(0, 1).
		Bold(true)
}

var statusRight = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("7")).
	Padding(0, 1)

var statusBar = lipgloss.NewStyle().
	Background(lipgloss.Color("0")).
	Height(1).
	Width(0)

var statusMid = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8")).
	Background(lipgloss.Color("0")).Align(lipgloss.Center)

var commandStyle = lipgloss.NewStyle()
