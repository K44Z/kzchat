package screens

import (
	"fmt"
	"kzchat/api"
	"kzchat/server/schemas"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type FocusArea int

const (
	ViewPort FocusArea = iota
	InputBox
	CommandBox
)

type ChatModel struct {
	Chat      schemas.Chat
	Messages  []schemas.Message
	Message   schemas.Message
	Command   string
	Current   schemas.User
	Textarea  textarea.Model
	Recipient schemas.User
	Err       string
	Channels  []string
	Width     int
	Height    int
	Ws        *websocket.Conn
	Viewport  viewport.Model
	Content   string
	Ready     bool
}

func NewChatModel(width int, height int) *ChatModel {
	// input := textinput.New()
	// input.Prompt = ""
	// input.Focus()
	// input.CharLimit = 500
	// input.Width = 65

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	// ta.Prompt = "┃ "
	ta.CharLimit = 500

	ta.SetWidth(100)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false
	vp := viewport.New(width, height-10)
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1

	currentUser := schemas.User{
		Username: api.Config.Username,
	}

	m := &ChatModel{
		Messages: []schemas.Message{},
		Current:  currentUser,
		Channels: []string{"general", "random", "dev", "help"},
		Width:    width,
		Height:   height,
		Textarea: ta,
		Err:      fmt.Sprint("width:", width),
		Viewport: vp,
	}
	str := m.RenderMessages()
	m.Viewport.SetContent(str)
	return m
}

func (m *ChatModel) Init() tea.Cmd {
	return m.ConnectToWs()
}

func (m *ChatModel) Update(msg tea.Msg, focusedArea int) (*ChatModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.Err = ""
		switch focusedArea {
		case 1:
			switch msg.String() {
			case "up", "k":
				m.Viewport.ScrollUp(1)
			case "down", "j":
				m.Viewport.ScrollDown(1)
			case "G":
				m.Viewport.GotoBottom()
			case "g":
				m.Viewport.GotoTop()
			default:
				m.Viewport.SetContent(m.RenderMessages())
				m.Viewport, cmd = m.Viewport.Update(msg)
				cmds = append(cmds, cmd)
			}
		case 2:
			m.Textarea, cmd = m.Textarea.Update(msg)
			cmds = append(cmds, cmd)
		}

	case ErrMsg:
		m.Err = msg.Error()

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Viewport.Width = msg.Width / 3
		m.Viewport.Height = msg.Height - 10
		m.Viewport.SetContent(m.RenderMessages())
	}

	return m, tea.Batch(cmds...)
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
		Border(lipgloss.ThickBorder(), true, true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#85c1dc")).
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

	textareaStyle := lipgloss.NewStyle().
		Width(m.Width).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#99d1db"))

	inputSection := textareaStyle.Render(m.Textarea.View())

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

func (m *ChatModel) renderLeftSidebar() string {
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

func (m *ChatModel) RenderMessages() string {
	var content strings.Builder
	if len(m.Messages) == 0 {
		content.WriteString(m.renderGuide(m.Width + 50))
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

	titleStyle := lipgloss.NewStyle().Align(lipgloss.Center).
		Foreground(lipgloss.Color("#8839ef"))

	content.WriteString(titleStyle.Render(fmt.Sprintf("%v", m.Chat.Name)))
	content.WriteString("\n\n")
	content.WriteString(headerStyle.Render("USER INFO") + "\n\n")
	content.WriteString(fmt.Sprintf("Name: %s\n", m.Current.Username))
	content.WriteString("Status: Online\n\n\n")
	content.WriteString(headerStyle.Render("SERVER") + "\n\n")
	content.WriteString("Connected\n")
	content.WriteString("Latency: 25ms\n")
	content.WriteString(fmt.Sprintf("\nMessages: %d", len(m.Messages)))
	return content.String()
}

func (m ChatModel) renderGuide(width int) string {
	content := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#8839ef")).
		Width(width + 10).
		Render(`
Welcome to kzchat. This chat app lets you talk to users in a clean, minimal interface.

Available commands:
  - open <username>         → Open chat with a user
  - dm <username> <message> → Send a direct message without opening chat

Controls:
  - Press [i]               → Enter input mode
  - Press [:]               → Open command bar
  - Press [esc]             → Exit current input mode
  - Press [q] or [Ctrl+C]   → Quit the app
`)

	return content
}
