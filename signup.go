package main

import (
	"bytes"
	"encoding/json"
	"kzchat/server/schemas"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "net/http"
)

var (
	primaryColor    = lipgloss.Color("#1e1e2e")
	secondaryColor  = lipgloss.Color("#00FFFF")
	backgroundColor = lipgloss.Color("#282C34")
	textColor       = lipgloss.Color("#FAFAFA")
	errorColor      = lipgloss.Color("#FF5252")
	mutedColor      = lipgloss.Color("#4c4f69")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Background(primaryColor).
			Padding(1, 4).
			MarginBottom(1).
			BorderStyle(lipgloss.RoundedBorder()).
			Align(lipgloss.Center)

	labelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			PaddingRight(1)

	inputBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1)

	focusedStyle = inputBoxStyle.Copy().
			BorderForeground(mutedColor)

	blurredStyle = inputBoxStyle.Copy().
			Foreground(mutedColor)

	cursorStyle = lipgloss.NewStyle().Foreground(secondaryColor)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	buttonStyle = lipgloss.NewStyle().
			Padding(0, 3).
			MarginTop(1).
			BorderStyle(lipgloss.RoundedBorder())

	focusedButton = buttonStyle.Copy().
			BorderForeground(secondaryColor).
			Foreground(secondaryColor).
			Bold(true).
			Render(" Submit ")

	blurredButton = buttonStyle.Copy().
			BorderForeground(mutedColor).
			Foreground(mutedColor).
			Render(" Submit ")
)

type SignupModel struct {
	focusIndex int
	inputs     []textinput.Model
	err        string
	message    string
}

func NewSignupModel() SignupModel {
	username := textinput.New()
	username.CharLimit = 256
	username.Focus()

	password := textinput.New()
	password.CharLimit = 256
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	return SignupModel{
		inputs:     []textinput.Model{username, password},
		focusIndex: 0,
	}
}

func (m SignupModel) Update(msg tea.Msg) (SignupModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			return m, func() tea.Msg {
				return screenMsg(loginScreen)
			}
		case "tab":
			m.focusIndex = (m.focusIndex + 1) % 3
			m.handleFocus()
		case "shift+tab":
			m.focusIndex = (m.focusIndex - 1) % 3
			m.handleFocus()
		case "enter":
			username := m.inputs[0].Value()
			password := m.inputs[1].Value()
			if username != "" && password != "" {
				auth := schemas.Auth{Username: username, Password: password}
				jsonData, err := json.Marshal(auth)
				if err != nil {
					m.err = "Error preparing request: " + err.Error()
					return m, nil
				}
				resp, err := http.Post("http://localhost:4000/auth/register", "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					m.err = "Error sending request: " + err.Error()
					return m, nil
				}

				defer resp.Body.Close()
				if resp.StatusCode == 201 {
					return m, func() tea.Msg {
						return screenMsg(loginScreen)
					}
				} else {
					var response map[string]string
					if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
						m.err = err.Error()
					} else {
						for i := range m.inputs {
							m.inputs[i].Reset()
							m.focusIndex = 0
							m.handleFocus()
						}
						m.err = response["message"]
					}
				}
			} else {
				m.err = "All Fields are required"
			}
		}

	}

	for i := range m.inputs {
		m.inputs[i], _ = m.inputs[i].Update(msg)
	}
	return m, nil
}

func (m SignupModel) View() string {
	var b strings.Builder

	b.WriteString("\n\n\n\n\n\n\n\n\n\n\n\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render("Username") + "\n")
	if m.focusIndex == 0 {
		b.WriteString(focusedStyle.Render(m.inputs[0].View()))
	} else {
		b.WriteString(blurredStyle.Render(m.inputs[0].View()))
	}
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render("Password") + "\n")
	if m.focusIndex == 1 {
		b.WriteString(focusedStyle.Render(m.inputs[1].View()))
	} else {
		b.WriteString(blurredStyle.Render(m.inputs[1].View()))
	}
	b.WriteString("\n")

	button := blurredButton
	if m.focusIndex == len(m.inputs) {
		button = focusedButton
	}
	b.WriteString(button + "\n\n")

	if m.err != "" {
		b.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Error: "+m.err))
	}

	return b.String()
}

func (m SignupModel) handleFocus() {
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}
