package  screens

import (
	"bytes"
	"encoding/json"
	"strings"

	authentication "kzchat/auth"
	"kzchat/server/schemas"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoginModel struct {
	inputs     []textinput.Model
	focusIndex int
	err        string
}

func NewLoginModel() *LoginModel {
	username := textinput.New()
	username.CharLimit = 256
	username.Focus()

	password := textinput.New()
	password.CharLimit = 256
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	return &LoginModel{
		inputs:     []textinput.Model{username, password},
		focusIndex: 0,
	}
}

func (m *LoginModel) Update(msg tea.Msg) (*LoginModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			return m, func() tea.Msg {
				return ScreenMsg(SignupScreen)
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
				resp, err := http.Post("http://localhost:4000/auth/login", "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					m.err = "Error sending request: " + err.Error()
					return m, nil
				}

				defer resp.Body.Close()
				var response map[string]string
				if resp.StatusCode == 200 {
					if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
						m.err = err.Error()
					} else {
						config := schemas.Config{Token: response["token"], Username: username}
						authentication.SaveConfig(config)
					}
					return m, func() tea.Msg {
						return ScreenMsg(ChatScreen)
					}
				} else {
					if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
						m.err = err.Error()
					} else {
						m.focusIndex = 0
						m.handleFocus()
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

func (m *LoginModel) View() string {

	var b strings.Builder

	b.WriteString("\n\n\n\n\n\n\n\n\n\n\n\n\n") // change this

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

func (m *LoginModel) handleFocus() {
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}
