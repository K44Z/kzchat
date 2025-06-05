package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

var quitKeys = key.NewBinding(
	key.WithKeys("ctrl+c", "ctrl+z"),
	key.WithHelp("", "press q to quit"),
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.chat.width = msg.Width 
		m.chat.height = msg.Height	
		return m, nil
	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		}
		switch {
		case !m.commandMode && msg.String() == "`":
			m.commandMode = true
			m.chat.commandMode = true
			m.command.Focus()

		case m.commandMode && msg.String() == "enter":
			cmd := m.command.Value()
			m.handleCommand(cmd) 
			m.command.Reset()
			m.command.Blur()
			m.commandMode = false
			m.chat.commandMode = false

		case (m.commandMode && msg.String() == "esc") || (m.commandMode && m.command.Value() == "") :
			m.command.Reset()
			m.command.Blur()
			m.commandMode = false
			m.chat.commandMode = false
		}
		switch m.currentScreen {
		case signupScreen:
			m.signup, cmd = m.signup.Update(msg)
			return m, cmd
		case loginScreen:
			m.login, cmd = m.login.Update(msg)
			return m, cmd
		case chatScreen:
			m.chat.width = m.width 
			m.chat.height = m.height
			m.chat, cmd = m.chat.Update(msg)
		}
	case screenMsg:
		m.currentScreen = screen(msg)
		// if m.currentScreen == chatScreen{
		// 	return m, fetchMessages(m.config, "")
		// }
		return m, nil
		// case dmMsg:
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	if m.commandMode {
		 var cmd tea.Cmd 
		 m.command, cmd = m.command.Update(msg)
		 cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) handleCommand(cmd string) {
	if m.currentScreen == chatScreen { 
		m.chat.command = cmd
	}
}

