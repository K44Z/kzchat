package main

import (
	authentication "kzchat/auth"
	"kzchat/helpers"
	"kzchat/server/schemas"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

type screen int
type screenMsg screen

const (
	signupScreen screen = iota
	loginScreen
	chatScreen
)

func NewModel() model {
	var m model
	config, err := authentication.ReadConfig()
	command := textinput.New()
	command.CharLimit = 256
	command.Prompt = ""
	m.command = command
	if config.Token == "" || err != nil || !authentication.IsTokenValid(config.Token){ // zid istokenvalid
		helpers.Logger.Println(err)
		m.config = config
		m.currentScreen = signupScreen
		m.login = NewLoginModel()
		m.signup = NewSignupModel()

	} else {
		m.currentScreen = chatScreen
		m.chat = NewChatModel(config.Username)
	}
	return m
}

type model struct {
	width         int
	height        int
	quitting      bool
	spinner       spinner.Model
	currentScreen screen
	tabs          []string
	signup        SignupModel
	login         LoginModel
	chat          ChatModel
	config        schemas.Config
	command       textinput.Model
	commandMode   bool
}
