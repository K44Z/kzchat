package main

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var content string
	var screenName string
	switch m.currentScreen {
	case loginScreen:
		content = m.login.View()
		screenName = "Login Screen"
	case signupScreen:
		content = m.signup.View()
		screenName = "Signup Screen"
	case chatScreen:
		content = m.chat.View()
		screenName = "Chat Screen"
	default:
		content = "Unknown screen"
		screenName = "Unknown"
	}
	if m.commandMode {
		screenName = "Command"
	}
	left := statusLeft.Render(screenName)
	right := statusRight.Render(" KZchat ")

	var change string
	if m.currentScreen == loginScreen {
		change = " • [`] signup"
	} else if m.currentScreen == signupScreen {
		change = " • [`] login"
	} else {
		change = ""
	}
	mid := statusMid.Render(" [tab] switch focus" + change + " • [ctrl+c] quit")

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
