package main

import (
	authentication "kzchat/auth"
	"kzchat/helpers"
	"kzchat/server/schemas"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int
type screenMsg screen

var config schemas.Config

const (
	signupScreen screen = iota
	loginScreen
	chatScreen
)

var err error

func Init() tea.Cmd {
	go func() {
		for msg := range messages {
			Program.Send(msg)
		}
	}()
	return nil
}

func NewModel() model {
	var m model
	config, err = authentication.ReadConfig()
	command := textinput.New()
	command.CharLimit = 256
	command.Prompt = ""
	m.command = command
	if config.Token == "" || err != nil || !authentication.IsTokenValid(config.Token) { // zid istokenvalid
		helpers.Logger.Println(err)
		m.currentScreen = signupScreen
		m.login = NewLoginModel()
		m.signup = NewSignupModel()
	} else {
		m.currentScreen = chatScreen
		m.chat = NewChatModel(m.width, m.height)
	}
	return m
}

type model struct {
	width         int
	height        int
	quitting      bool
	spinner       spinner.Model
	currentScreen screen
	signup        SignupModel
	login         LoginModel
	chat          ChatModel
	command       textinput.Model
	commandMode   bool
}
