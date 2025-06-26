package main

import (
	"fmt"
	repository "kzchat/server/database/generated"
	"kzchat/server/schemas"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

// var sampleMessages = []schemas.Message{
// 	{
// 		Content:        "Hey, are you there?",
// 		Time:           time.Now(),
// 		SenderUsername: "Alice",
// 	},
// 	{
// 		Content:        "Yeah, what's up?",
// 		Time:           time.Now(),
// 		SenderUsername: "Bob",
// 	},
// 	{
// 		Content:        "Just testing this Bubble Tea UI 😄",
// 		Time:           time.Date(2023, 1, 1, 10, 3, 0, 0, time.Local),
// 		SenderUsername: "Alice",
// 	},
// 	{
// 		Content:        "Oh nice! I love terminal UIs.",
// 		Time:           time.Date(2023, 1, 1, 10, 4, 0, 0, time.Local),
// 		SenderUsername: "Bob",
// 	},
// 	{
// 		Content:        "Same here. Super fun to build.",
// 		Time:           time.Date(2023, 1, 1, 10, 5, 0, 0, time.Local),
// 		SenderUsername: "Alice",
// 	},
// 	{
// 		Content:        "Do you want to add input next?",
// 		Time:           time.Date(2023, 1, 1, 10, 6, 0, 0, time.Local),
// 		SenderUsername: "Bob",
// 	},
// 	{
// 		Content:        "Yup! Let's do it.",
// 		Time:           time.Date(2023, 1, 1, 10, 7, 0, 0, time.Local),
// 		SenderUsername: "Alice",
// 	},
// 	{
// 		Content:        "This is a simulated system message.",
// 		Time:           time.Date(2023, 1, 1, 10, 8, 0, 0, time.Local),
// 		SenderUsername: "System",
// 	},
// 	{
// 		Content:        "Okay, now I'm just spamming.",
// 		Time:           time.Date(2023, 1, 1, 10, 9, 0, 0, time.Local),
// 		SenderUsername: "Bob",
// 	},
// 	{
// 		Content:        "😂",
// 		Time:           time.Date(2023, 1, 1, 10, 10, 0, 0, time.Local),
// 		SenderUsername: "Alice",
// 	},
// }

type ChatModel struct {
	chat        schemas.Chat
	messages    []schemas.Message
	message     repository.Message
	commandMode bool
	command     string
	input       textinput.Model
	username    string
	recipient   string
	err         string
	channels    []string
	width       int
	height      int
	ws          *websocket.Conn
	viewport    viewport.Model
	content     string
	ready       bool
}

func NewChatModel(width int, height int) ChatModel {
	input := textinput.New()
	input.Prompt = ""
	input.Focus()
	input.CharLimit = 500
	input.Width = 65

	vp := viewport.New(width/3, height-10)
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1

	m := ChatModel{
		messages: []schemas.Message{},
		username: config.Username,
		input:    input,
		channels: []string{"general", "random", "dev", "help"},
		width:    width,
		height:   height,
		viewport: vp,
	}
	str := m.renderMessages()
	m.viewport.SetContent(str)
	return m
}

func (m ChatModel) Init() tea.Cmd {
	return m.connectToWs()
}

func (m ChatModel) Update(msg tea.Msg) (ChatModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			inputValue := strings.TrimSpace(m.input.Value())
			if inputValue != "" {
				message := schemas.Message{
					Content:        inputValue,
					SenderUsername: config.Username,
					Time:           time.Now(),
					Chat:           m.chat,
					ReceiverUsername: "username",
				}
				if m.ws == nil {
					if m.err != "" {
						m.err += "; ws is not connected"
					} else {
						m.err = "ws is not connected"
					}
					return m, nil
				}
				err := m.ws.WriteJSON(message)
				if err != nil {
					m.err = err.Error()
				}
				m.messages = append(m.messages, message)
				m.viewport.SetContent(m.renderMessages())
				m.viewport.GotoBottom()
				m.input.Reset()
			}
		case "up", "k":
			m.viewport.ScrollUp(1)
		case "down", "j":
			m.viewport.ScrollDown(1)
		}
	}
	if !m.commandMode {
		m.input, cmd = m.input.Update(msg)
		m.viewport, _ = m.viewport.Update(msg)
	}
	return m, cmd
}

func (m ChatModel) View() string {

	leftWidth := int(float64(m.width) * 0.14)
	rightWidth := int(float64(m.width) * 0.14)
	chatWidth := m.width - leftWidth - rightWidth - 7
	contentHeight := m.height - 7

	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(contentHeight).
		Padding(1, 1).
		Align(lipgloss.Left)

	chatStyle := lipgloss.NewStyle().
		Width(chatWidth).
		Height(contentHeight).
		Padding(1, 1).
		Align(lipgloss.Left).
		Border(lipgloss.NormalBorder(), false, true, false, true).
		BorderForeground(lipgloss.Color("240"))

	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(contentHeight).
		Padding(1, 1).
		Align(lipgloss.Left)

	widthString := fmt.Sprint(m.width)
	heightString := fmt.Sprint(m.height)

	leftSection := leftStyle.Render(m.renderLeftSidebar())
	chatSection := chatStyle.Render(m.viewport.View())
	rightSection := rightStyle.Render(m.renderRightSidebar() + widthString + " " + heightString)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, chatSection, rightSection)

	inputStyle := lipgloss.NewStyle().
		Width(m.width-4).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("240"))

	inputSection := inputStyle.Render(m.input.View())

	var content strings.Builder
	content.WriteString(mainContent)
	content.WriteString(inputSection)

	if m.err != "" {
		errorStyle := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("9")).
			Bold(true)

		errorMsg := errorStyle.Render("Error: " + m.err)
		content.WriteString("\n\n" + errorMsg)
	}
	return content.String()
}

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
		// messageStyle := lipgloss.NewStyle().
		// 	Foreground(lipgloss.Color("15"))

		// timestampStyle := lipgloss.NewStyle().
		// 	Foreground(lipgloss.Color("240"))

		// usernameStyle := lipgloss.NewStyle().
		// 	Bold(true).
		// 	Foreground(lipgloss.Color("14"))

		// for _, msg := range sampleMessages {
		// 	parts := strings.SplitN(msg.Content, "] ", 2)
		// 	if len(parts) == 2 {
		// 		timestamp := timestampStyle.Render(parts[0] + "]")

		// 		msgParts := strings.SplitN(parts[1], ": ", 2)
		// 		if len(msgParts) == 2 {
		// 			username := usernameStyle.Render(msgParts[0])
		// 			message := messageStyle.Render(msgParts[1])
		// 			content.WriteString(fmt.Sprintf("%s %s: %s\n", timestamp, username, message))
		// 		}
		// 	}
		// }

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
			timestamp := timestampStyle.Render(fmt.Sprintf("[%s]", msg.Time.Format("15:04")))
			username := usernameStyle.Render(string(msg.SenderUsername))
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
	content.WriteString("Latency: 25ms\n")
	content.WriteString(fmt.Sprintf("\nMessages: %d", len(m.messages)))
	return content.String()
}
