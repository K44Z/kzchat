package main

import (
	s "github.com/K44Z/kzchat/pkg/screens"

	"github.com/K44Z/kzchat/internal/api"
	"github.com/K44Z/kzchat/internal/helpers"

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
	List
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

	width, height := m.width, m.height
	if width == 0 {
		width = 20
	}
	if height == 0 {
		height = 30
	}

	items, err := api.GetUsers()
	if err != nil {
		helpers.Logger.Printf("Error fetching users: %v", err)
	}

	userList := list.New(items, s.UserDelegate{}, width, height)
	userList.Title = "Users"
	userList.SetShowTitle(true)
	userList.SetShowFilter(true)
	userList.SetFilteringEnabled(true)
	userList.SetShowHelp(true)
	userList.SetShowStatusBar(true)
	defaultStyles := list.DefaultStyles()
	userList.Styles.Title = defaultStyles.Title
	userList.Styles.FilterCursor = defaultStyles.FilterCursor
	userList.Styles.FilterPrompt = defaultStyles.FilterPrompt
	m.List = userList

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
	currentScreen api.Screen
	signup        *s.SignupModel
	login         *s.LoginModel
	chat          *s.ChatModel
	command       textinput.Model
	FocusArea     FocusArea
	List          list.Model
}
