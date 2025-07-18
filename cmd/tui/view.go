package main

import (
	"math"
	"strings"

	"github.com/K44Z/kzchat/pkg/screens"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var content string
	var mode string
	switch m.currentScreen {
	case screens.LoginScreen:
		content = m.login.View()
		mode = "Login Screen"
	case screens.SignupScreen:
		content = m.signup.View()
		mode = "Signup Screen"
	case screens.ChatScreen:
		content = m.chat.View()
		switch m.FocusArea {
		case 1:
			mode = "Normal"
		case 3:
			mode = "Command"
		case 2:
			mode = "Insert"
		default:
			mode = "Chat Screen"
		}
	default:
		content = "Unknown screen"
		mode = "Unknown"
	}

	left := StatusStyle(mode).Render(mode)
	right := StatusStyle(mode).Render((" KZchat "))

	var change string

	switch m.currentScreen {
	case screens.LoginScreen:
		change = " • [`] signup"
	case screens.SignupScreen:
		change = " • [`] login"
	default:
		change = ""
	}
	mid := statusMid.Render(change + " • [ctrl+z] quit")

	gap := int(math.Max(0, float64(m.width-lipgloss.Width(left+right+mid))))
	gapString := statusMid.Render(strings.Repeat(" ", gap))
	bar := statusBar.Width(m.width).MarginBottom(0).Render(left + mid + gapString + right)
	command := commandStyle.Render(m.command.View())
	box := layoutStyle.
		Width(m.width).
		Height(m.height - 2).
		Render(content)

	return box + "\n" + bar + "\n" + command
}
