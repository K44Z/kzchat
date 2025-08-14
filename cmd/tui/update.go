package main

import (
	"fmt"
	"strings"

	"github.com/K44Z/kzchat/internal/api"
	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/K44Z/kzchat/pkg/screens"
	s "github.com/K44Z/kzchat/pkg/screens"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	return nil
}

var quitKeys = key.NewBinding(
	key.WithKeys("ctrl+z"),
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		switch m.currentScreen {
		case screens.ChatScreen:
			var cmd tea.Cmd
			m.chat, cmd = m.chat.Update(msg, int(m.FocusArea))
			cmds = append(cmds, cmd)
		case screens.LoginScreen:
			if m.login != nil {
				var cmd tea.Cmd
				m.login, cmd = m.login.Update(msg)
				cmds = append(cmds, cmd)
			}
		case screens.SignupScreen:
			if m.signup != nil {
				var cmd tea.Cmd
				m.signup, cmd = m.signup.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
		m.command.Width = msg.Width - 4
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		}

		switch msg.String() {
		case ":":
			if m.FocusArea == 1 {
				var cmd tea.Cmd
				m.FocusArea = 3
				m.command.Focus()
				m.command, cmd = m.command.Update(msg)
				return m, cmd
			}
		case "i":
			if m.currentScreen == s.ChatScreen {
				if m.FocusArea == 1 {
					m.FocusArea = 2
					m.chat.Textarea.Focus()
					return m, nil
				}
			}
		case "esc":
			var cmd tea.Cmd
			switch m.FocusArea {
			case 3:
				m.command.Reset()
				m.command.Blur()
				m.FocusArea = 1
			case 2:
				m.chat.Textarea.Blur()
				m.FocusArea = 1
			case 4:
				if !m.List.FilterInput.Focused() {
					m.FocusArea = 1
					m.chat, cmd = m.chat.Update(msg, 1)
					return m, cmd
				}
			}
			m.FocusArea = 1
			return m, cmd
		case "/":
			if m.FocusArea == 4 {
				m.List.FilterInput.Focus()
				return m, nil
			}
		case "enter":
			switch m.FocusArea {
			case 3:
				command := strings.TrimSpace(m.command.Value())
				var cmd tea.Cmd
				if command != "" {
					cmd = m.handleCommand(command)
				}
				m.command.Reset()
				m.command.Blur()
				m.FocusArea = 2
				m.chat.Textarea.Focus()
				return m, cmd
			case 2:
				m.chat.SendMessage()
				m.chat.Textarea.Reset()
				m.chat.Textarea.Focus()
			case 4:
				var cmd tea.Cmd
				selectedItem := m.List.SelectedItem().(schemas.User)
				cmd, err = m.handleRecipientSelection(selectedItem)
				if err != nil {
					return m, nil
				}
				m.FocusArea = 1
				return m, cmd
			}
		case "ctrl+s":
			m.FocusArea = 4
		}

		switch m.FocusArea {
		case 3:
			var cmd tea.Cmd
			m.command, cmd = m.command.Update(msg)
			cmds = append(cmds, cmd)
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
			m.chat, cmd = m.chat.Update(msg, int(m.FocusArea))
			cmds = append(cmds, cmd)
		}

	case api.ScreenMsg:
		m.currentScreen = api.Screen(msg)
		if m.currentScreen == screens.ChatScreen {
			m.chat = screens.NewChatModel(m.width, m.height)
			cmd := m.chat.Init()
			m.SetList()
			m.FocusArea = 2
			return m, cmd
		}
		return m, nil

	case api.WsMsg:
		scMesasge := schemas.Message{
			Content: msg.Content,
			Time:    msg.Time,
			Sender:  msg.Sender,
		}
		m.chat.Messages = append(m.chat.Messages, scMesasge)
		var cmd tea.Cmd
		m.chat, cmd = m.chat.Update(msg, int(m.FocusArea))
		return m, cmd

	case api.WsConnectedMsg:
		m.chat.Ws = msg.Conn
		go m.chat.ReadLoop(msg.Conn)
		return m, nil

	case api.ChatFetchedMsg:
		m.chat.Messages = msg.Messages
		m.chat.Viewport.SetContent(m.chat.RenderMessages())
		m.chat.Viewport.GotoBottom()
		return m, nil

	case api.ErrMsg:
		m.chat.Err = msg.Error()

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *model) handleCommand(c string) tea.Cmd {
	if m.currentScreen == screens.ChatScreen {
		m.chat.Command = c
		cmd := m.chat.HandleChatCommand()
		return cmd
	}
	return nil
}

func (m *model) handleRecipientSelection(r schemas.User) (tea.Cmd, error) {
	id, users, err := api.GetChat([]string{api.Config.Username, r.Username})
	if users == nil || err != nil {
		return nil, err
	}
	chat := schemas.Chat{
		Name: fmt.Sprint(api.Config.Username, " - ", r.Username),
		ID:   *id,
	}
	m.chat.Chat = chat
	m.chat.Recipient.Username = r.Username
	cmd := m.chat.FetchMessages()
	m.chat.Current = users[0]
	m.chat.Recipient = users[1]
	return cmd, nil
}

func (m *model) SetList() {
	width, height := m.width, m.height
	if width == 0 {
		width = 20
	}
	if height == 0 {
		height = 30
	}
	items, err := api.GetUsers()
	if err != nil {
		m.chat.Err = "Error fetching the list of users"
	}

	delegate := list.NewDefaultDelegate()

	userList := list.New(items, delegate, width, height)
	userList.Title = "Users"

	m.List = userList
}
