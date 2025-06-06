package main

import (
	"fmt"
	"strings"

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

var (
	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

var layoutStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Margin(0)

var statusLeft = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#eff1f5")).
	Background(lipgloss.Color("#8839ef")).
	Padding(0, 1).
	Bold(true).
	BorderRight(true).
	BorderForeground(lipgloss.Color("#89b4fa"))

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

func (m ChatModel) renderLeftSidebar() string {
	var content strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#8839ef"))
	content.WriteString(headerStyle.Render("Channels") + "\n\n")

	if len(m.channels) == 0 {
		content.WriteString("No channels")
	} else {
		for i, ch := range m.channels {
			channelStyle := lipgloss.NewStyle()
			prefix := " "

			if i == 0 {
				channelStyle = channelStyle.
					Bold(true).
					Foreground(lipgloss.Color("15"))
				prefix = ">"
			}

			channel := channelStyle.Render(fmt.Sprintf("%s #%s", prefix, ch))
			content.WriteString(channel + "\n")
		}
	}

	content.WriteString("\n\n")

	return content.String()
}

func (m ChatModel) renderMessages() string {
	var content strings.Builder

	if len(m.messages) == 0 {
		messageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

		timestampStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

		usernameStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14"))

		messages := []string{
			"[12:34] system: Welcome to KZ Chat",
			"[12:35] admin: Server is online",
			"[12:36] user1: Hello everyone",
		}

		for _, msg := range messages {
			parts := strings.SplitN(msg, "] ", 2)
			if len(parts) == 2 {
				timestamp := timestampStyle.Render(parts[0] + "]")

				msgParts := strings.SplitN(parts[1], ": ", 2)
				if len(msgParts) == 2 {
					username := usernameStyle.Render(msgParts[0])
					message := messageStyle.Render(msgParts[1])
					content.WriteString(fmt.Sprintf("%s %s: %s\n", timestamp, username, message))
				}
			}
		}

		content.WriteString("\n")
		promptStyle := lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("240"))
		content.WriteString(promptStyle.Render("Type a message to start chatting..."))
	} else {
		timestampStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

		usernameStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14"))

		messageStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

		for _, msg := range m.messages {
			timestamp := timestampStyle.Render(fmt.Sprintf("[%s]", msg.Time.Time.Format("15:04")))
			username := usernameStyle.Render(string(msg.SenderID))
			message := messageStyle.Render(msg.Content)

			content.WriteString(fmt.Sprintf("%s %s: %s\n", timestamp, username, message))
		}
	}

	return content.String()
}

func (m ChatModel) renderRightSidebar() string {
	var content strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true)
	content.WriteString(headerStyle.Render("USER INFO") + "\n\n")
	content.WriteString(fmt.Sprintf("Name: %s\n", m.username))
	content.WriteString("Status: Online\n\n\n")
	content.WriteString(headerStyle.Render("SERVER") + "\n\n")
	content.WriteString("Connected\n")
	content.WriteString("Latency: 25ms\n\n\n")
	return content.String()
}
