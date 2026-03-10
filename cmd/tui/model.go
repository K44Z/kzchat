package main

import (
	s "github.com/K44Z/kzchat/pkg/screens"

	"github.com/K44Z/kzchat/internal/api"
	"github.com/K44Z/kzchat/internal/messages"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var err error

type FocusArea int

const (
	ViewPort FocusArea = iota
	InputBox
	CommandBox
	UserList
	FilePicker
	AttachmentList
)

func (m *model) Init() tea.Cmd {
	return func() tea.Msg {
		configErr := api.ReadConfig()
		if api.Config.Token == "" || configErr != nil {
			return messages.ScreenMsg(s.LoginScreen)
		}
		err := api.IsTokenValid(api.Config.Token)
		if err != nil {
			m.login.Err = err.Error()
			return messages.ScreenMsg(s.LoginScreen)
		}
		return messages.ScreenMsg(s.ChatScreen)
	}
}

func NewModel() model {
	var m model
	command := textinput.New()
	command.CharLimit = 256
	command.Prompt = ""
	m.command = command
	m.login = s.NewLoginModel()
	m.signup = s.NewSignupModel()
	return m
}

type model struct {
	width          int
	height         int
	quitting       bool
	spinner        spinner.Model
	currentScreen  messages.Screen
	signup         *s.SignupModel
	login          *s.LoginModel
	chat           *s.ChatModel
	command        textinput.Model
	FocusArea      FocusArea
	List           list.Model
	AttachmentList list.Model
}
