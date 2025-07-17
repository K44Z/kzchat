package screens

import (
	"fmt"
	"kzchat/api"
	"kzchat/helpers"
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
	Chat          schemas.Chat
	Messages      []schemas.Message
	Message       schemas.Message
	Command       string
	Current       schemas.User
	Textarea      textarea.Model
	Recipient     schemas.User
	Err           string
	Channels      []string
	Width         int
	ChatWidth     int
	LeftWidth     int
	RightWidth    int
	Height        int
	ContentHeight int
	Ws            *websocket.Conn
	Viewport      viewport.Model
	Content       string
	Ready         bool
}

func NewChatModel(width int, height int) *ChatModel {
	leftWidth := helpers.ComputeSideWidth(width)
	rightWidth := helpers.ComputeSideWidth(width)
	chatWidth := helpers.ComputeChatWidth(width, leftWidth, rightWidth)
	contentHeight := helpers.ComputeContentHeight(height)

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.CharLimit = 500

	textareaWidth := chatWidth - 4
	if textareaWidth < 20 {
		textareaWidth = 20
	}

	ta.SetWidth(textareaWidth)
	ta.SetHeight(1)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	viewportWidth := chatWidth - 4
	viewportHeight := contentHeight - 6
	if viewportWidth < 20 {
		viewportWidth = 20
	}
	if viewportHeight < 5 {
		viewportHeight = 5
	}

	vp := viewport.New(viewportWidth, viewportHeight)
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 1

	currentUser := schemas.User{
		Username: api.Config.Username,
	}

	m := &ChatModel{
		Messages:      []schemas.Message{},
		Current:       currentUser,
		Channels:      []string{"general", "random", "dev", "help"},
		Width:         width,
		Height:        height,
		Textarea:      ta,
		Err:           "",
		Viewport:      vp,
		ContentHeight: contentHeight,
		ChatWidth:     chatWidth,
		RightWidth:    rightWidth,
		LeftWidth:     leftWidth,
		Ready:         true,
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
		// m.Err = fmt.Sprintf("width : %v, height : %v", m.Width, m.Height)
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

		m.LeftWidth = helpers.ComputeSideWidth(m.Width)
		m.RightWidth = helpers.ComputeSideWidth(m.Width)
		m.ChatWidth = helpers.ComputeChatWidth(m.Width, m.LeftWidth, m.RightWidth)
		m.ContentHeight = helpers.ComputeContentHeight(m.Height)

		m.Textarea.SetWidth(m.ChatWidth - 4)

		viewportWidth := m.ChatWidth - 4
		viewportHeight := m.ContentHeight - 6
		if viewportWidth < 20 {
			viewportWidth = 20
		}
		if viewportHeight < 5 {
			viewportHeight = 5
		}

		oldYPosition := m.Viewport.YPosition
		m.Viewport = viewport.New(viewportWidth, viewportHeight)
		m.Viewport.MouseWheelEnabled = true
		m.Viewport.MouseWheelDelta = 1
		m.Viewport.SetContent(m.RenderMessages())

		m.Viewport.YPosition = oldYPosition
		m.Ready = true
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
	if !m.Ready {
		return "loading ..."
	}

	compactMode := m.Width < 100

	contentBoxHeight := m.ContentHeight - 2
	if contentBoxHeight < 5 {
		contentBoxHeight = 5
	}

	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), true, true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#414559")).
		Width(m.LeftWidth).
		Height(contentBoxHeight).
		Padding(1, 1).
		Align(lipgloss.Left)

	chatBoxWidth := m.ChatWidth - 8
	if chatBoxWidth < 20 {
		chatBoxWidth = 20
	}

	chatStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), true, true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#414559")).
		Width(chatBoxWidth).
		Height(contentBoxHeight).
		Padding(1, 1).
		Align(lipgloss.Left)

	rightStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), true, true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#414559")).
		Width(m.RightWidth).
		Height(contentBoxHeight).
		Padding(1, 1).
		Align(lipgloss.Left)

	leftSection := leftStyle.Render(m.renderLeftSidebar())
	chatSection := chatStyle.Render(m.Viewport.View())
	rightSection := rightStyle.Render(m.renderRightSidebar())

	var mainContent string
	if compactMode {
		mainContent = chatSection
	} else {
		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, leftSection, chatSection, rightSection)
	}
	var textareaWidth int
	if compactMode {
		textareaWidth = chatBoxWidth
	} else {
		textareaWidth = m.Width - 10
	}

	if textareaWidth < 20 {
		textareaWidth = 20
	}

	textareaStyle := lipgloss.NewStyle().
		Width(textareaWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#414559"))

	inputSection := textareaStyle.Render(m.Textarea.View())

	var content strings.Builder
	content.WriteString(mainContent)
	content.WriteString("\n" + inputSection)

	if m.Err != "" {
		errorStyle := lipgloss.NewStyle().
			Width(m.Width).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("9")).
			Bold(true)

		errorMsg := errorStyle.Render("Error: " + m.Err)
		content.WriteString("\n" + errorMsg)
	}
	return content.String()
}

func (m *ChatModel) renderLeftSidebar() string {
	var content strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#8839ef"))

	content.WriteString(headerStyle.Render("Channels"))
	content.WriteString("\n\n")

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
	messageAreaWidth := m.ChatWidth - 10
	if messageAreaWidth < 20 {
		messageAreaWidth = 20
	}

	if len(m.Messages) == 0 {
		content.WriteString(m.renderGuide(messageAreaWidth))
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

func (m *ChatModel) renderRightSidebar() string {
	var content strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true)

	titleStyle := lipgloss.NewStyle().Align(lipgloss.Center).
		Foreground(lipgloss.Color("#8839ef"))

	content.WriteString(titleStyle.Render(fmt.Sprintf("%v", m.Chat.Name)))
	content.WriteString("\n\n")
	content.WriteString(headerStyle.Render("USER INFO"))
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("Name: %s\n", m.Current.Username))
	content.WriteString("Status: Online\n\n\n")
	content.WriteString(headerStyle.Render("SERVER"))
	content.WriteString("\n\n")
	content.WriteString("Connected\n")
	content.WriteString("Latency: 25ms\n")

	return content.String()
}

func (m *ChatModel) renderGuide(width int) string {
	guideText := `
Welcome to kzchat. This chat app lets you talk to users in a clean, minimal interface.

Available commands:
  - open <username>         → Open chat with a user
  - dm <username> <message> → Send a direct message without opening chat

Controls:
  - Press [i]               → Enter input mode
  - Press [:]               → Open command bar
  - Press [esc]             → Exit current input mode
  - Press [Ctrl+z]          → Quit the app
`

	renderWidth := min(width, 100)
	content := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#8839ef")).
		Width(renderWidth).
		Render(guideText)

	return content
}
