package screens

import (
	"fmt"
	"kzchat/api"
	"kzchat/server/schemas"
	"strings"

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
	Chat        schemas.Chat
	Messages    []schemas.Message
	Message     schemas.Message
	CommandMode bool
	Command     string
	Input       textinput.Model
	Current     schemas.User
	Recipient   schemas.User
	Err         string
	Channels    []string
	Width       int
	Height      int
	Ws          *websocket.Conn
	Viewport    viewport.Model
	Content     string
	Ready       bool
}

func NewChatModel(width int, height int) *ChatModel {
	input := textinput.New()
	input.Prompt = ""
	input.Focus()
	input.CharLimit = 500
	input.Width = 65

	vp := viewport.New(width/3, height-10)
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1

	currentUser := schemas.User{
		Username: api.Config.Username,
	}

	m := &ChatModel{
		Messages: []schemas.Message{},
		Current:  currentUser,
		Input:    input,
		Channels: []string{"general", "random", "dev", "help"},
		Width:    width,
		Height:   height,
		Viewport: vp,
	}
	str := m.RenderMessages()
	m.Viewport.SetContent(str)
	return m
}

func (m *ChatModel) Init() tea.Cmd {
	return m.ConnectToWs()
}

func (m *ChatModel) Update(msg tea.Msg) (*ChatModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.Err = ""
		switch msg.String() {
		case "enter":
			if m.CommandMode && msg.String() == "enter" {
				cmd := m.HandleChatCommand()
				m.CommandMode = false
				return m, cmd
			} else {
				m.SendMessage()
				return m, nil
			}
		case "up", "k":
			m.Viewport.ScrollUp(1)
		case "down", "j":
			m.Viewport.ScrollDown(1)
		}
	}
	if !m.CommandMode {
		m.Input, cmd = m.Input.Update(msg)
		m.Viewport, _ = m.Viewport.Update(msg)
	}
	return m, cmd
}

func (m *ChatModel) HandleChatCommand() tea.Cmd {
	var (
		cmd    tea.Cmd
		output string
	)
	parts := strings.Fields(m.Command)
	if len(parts) == 0 {
		return nil
	}

	name := parts[0][1:]
	args := parts[1:]
	if handler, ok := CommandRegistry[name]; ok {
		output, cmd = handler(CommandContext{Model: m}, args)
		if output != "" {
			return func() tea.Msg {
				return ErrMsg(fmt.Errorf("%s", output))
			}
		}
	} else {
		m.Err = "Unknown command: " + name
	}
	return cmd
}

func (m *ChatModel) View() string {

	leftWidth := int(float64(m.Width) * 0.14)
	rightWidth := int(float64(m.Width) * 0.14)
	chatWidth := m.Width - leftWidth - rightWidth - 7
	contentHeight := m.Height - 7

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

	widthString := fmt.Sprint(m.Width)
	heightString := fmt.Sprint(m.Height)

	leftSection := leftStyle.Render(m.renderLeftSidebar())
	chatSection := chatStyle.Render(m.Viewport.View())
	rightSection := rightStyle.Render(m.renderRightSidebar() + widthString + " " + heightString)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, chatSection, rightSection)

	inputStyle := lipgloss.NewStyle().
		Width(m.Width-4).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("240"))

	inputSection := inputStyle.Render(m.Input.View())

	var content strings.Builder
	content.WriteString(mainContent)
	content.WriteString(inputSection)

	if m.Err != "" {
		errorStyle := lipgloss.NewStyle().
			Width(m.Width).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("9")).
			Bold(true)

		errorMsg := errorStyle.Render("Error: " + m.Err)
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

	if len(m.Channels) == 0 {
		content.WriteString("No channels")
	} else {
		for i, ch := range m.Channels {
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

func (m ChatModel) RenderMessages() string {
	var content strings.Builder

	if len(m.Messages) == 0 {
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

		for _, msg := range m.Messages {
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
	content.WriteString(fmt.Sprintf("Name: %s\n", m.Current.Username))
	content.WriteString("Status: Online\n\n\n")
	content.WriteString(headerStyle.Render("SERVER") + "\n\n")
	content.WriteString("Connected\n")
	content.WriteString("Latency: 25ms\n")
	content.WriteString(fmt.Sprintf("\nMessages: %d", len(m.Messages)))
	return content.String()
}
