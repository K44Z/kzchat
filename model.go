package main

import (
	a "kzchat/api"
	"kzchat/helpers"
	s "kzchat/screens"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var err error

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
	a.ReadConfig()
	command := textinput.New()
	command.CharLimit = 256
	command.Prompt = ""
	m.command = command
	if a.Config.Token == "" || err != nil || !a.IsTokenValid(a.Config.Token) { // zid istokenvalid
		helpers.Logger.Println(err)
		m.currentScreen = s.SignupScreen
		m.login = s.NewLoginModel()
		m.signup = s.NewSignupModel()
	} else {
		m.currentScreen = s.ChatScreen
		m.chat = s.NewChatModel(m.width, m.height)
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
	commandMode   bool
}
