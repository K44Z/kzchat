package main

import (
	"kzchat/screens"
	"kzchat/server/schemas"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	return nil
}

var quitKeys = key.NewBinding(
	key.WithKeys("ctrl+c", "ctrl+z", "q"),
	key.WithHelp("q", "quit"),
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.chat != nil {
			m.chat.Width = msg.Width
			m.chat.Height = msg.Height
		}
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		}

		switch {
		case !m.commandMode && msg.String() == "`":
			m.commandMode = true
			m.chat.CommandMode = true
			m.command.Focus()

		case m.commandMode && msg.String() == "enter":
			cmd := m.command.Value()
			m.handleCommand(cmd)
			m.command.Reset()
			m.command.Blur()
			m.commandMode = false

		case m.commandMode && (msg.String() == "esc" || m.command.Value() == ""):
			m.command.Reset()
			m.command.Blur()
			m.commandMode = false
			m.chat.CommandMode = false
		}

		switch m.currentScreen {
		case screens.SignupScreen:
			var cmd tea.Cmd
			m.signup, cmd = m.signup.Update(msg)
			return m, cmd

		case screens.LoginScreen:
			var cmd tea.Cmd
			m.login, cmd = m.login.Update(msg)
			return m, cmd

		case screens.ChatScreen:
			m.chat.Width = m.width
			m.chat.Height = m.height
			var cmd tea.Cmd
			m.chat, cmd = m.chat.Update(msg)
			cmds = append(cmds, cmd)
		}
	case screens.ScreenMsg:
		m.currentScreen = screens.Screen(msg)
		if m.currentScreen == screens.ChatScreen {
			m.chat = screens.NewChatModel(m.width, m.height)
			cmd := m.chat.Init()
			// cmds = append(cmds, fetchMessages(m.chat.recipient))
			return m, cmd
		}
		return m, nil

	case screens.WsMsg:
		scMesasge := schemas.Message{
			Content:        msg.Content,
			Time:           msg.Time,
			SenderUsername: msg.SenderUsername,
		}
		m.chat.Messages = append(m.chat.Messages, scMesasge)
		var cmd tea.Cmd
		m.chat, cmd = m.chat.Update(msg)
		return m, cmd

	case screens.WsConnectedMsg:
		m.chat.Ws = msg.Conn
		go m.chat.ReadLoop(msg.Conn)
		return m, nil
	case screens.ChatFetchedMsg:
		m.chat.Messages = msg.Messages
		m.chat.Viewport.SetContent(m.chat.RenderMessages())
		m.chat.Viewport.GotoBottom()
		return m, nil
	case screens.ErrMsg:
		m.chat.Err = msg.Error()
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.commandMode {
		var cmd tea.Cmd
		m.command, cmd = m.command.Update(msg)
		cmds = append(cmds, cmd)
	}else {
		m.chat.Input.Focus()
	}
	return m, tea.Batch(cmds...)
}

func (m *model) handleCommand(cmd string) {
	if m.currentScreen == screens.ChatScreen {
		m.chat.Command = cmd
	}
}
