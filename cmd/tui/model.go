package main

import (
	s "github.com/K44Z/kzchat/pkg/screens"

	"github.com/K44Z/kzchat/internal/api"
	"github.com/K44Z/kzchat/internal/helpers"

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
)

func Init() tea.Cmd {
	go func() {
		for msg := range s.Messages {
			Program.Send(msg)
		}
	}()
	return nil
}

func NewModel() model {
	var m model
	api.ReadConfig()
	command := textinput.New()
	command.CharLimit = 256
	command.Prompt = ""
	m.command = command
	if api.Config.Token == "" || err != nil || !api.IsTokenValid(api.Config.Token) {
		helpers.Logger.Println(err)
		m.currentScreen = s.LoginScreen
		m.login = s.NewLoginModel()
		m.signup = s.NewSignupModel()
	} else {
		m.currentScreen = s.ChatScreen
		m.chat = s.NewChatModel(m.width, m.height)
		m.FocusArea = 1
	}
	return m
}

type model struct {
	width         int
	height        int
	quitting      bool
	spinner       spinner.Model
	currentScreen s.Screen
	signup        *s.SignupModel
	login         *s.LoginModel
	chat          *s.ChatModel
	command       textinput.Model
	FocusArea     FocusArea
}
