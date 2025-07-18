package screens

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor   = lipgloss.Color("#1e1e2e")
	secondaryColor = lipgloss.Color("#00FFFF")
	mutedColor     = lipgloss.Color("#4c4f69")
	inputBoxStyle  = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1)

	focusedStyle = inputBoxStyle.Copy().
			BorderForeground(mutedColor)

	blurredStyle = inputBoxStyle.Copy().
			Foreground(mutedColor)
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
}

func NewSignupModel() *SignupModel {
	username := textinput.New()
	username.CharLimit = 256
	username.Focus()

	password := textinput.New()
	password.CharLimit = 256
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	return &SignupModel{
		inputs:     []textinput.Model{username, password},
		focusIndex: 0,
	}
}

func (m *SignupModel) Update(msg tea.Msg) (*SignupModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.err = ""
		switch msg.String() {
		case "`":
			return m, func() tea.Msg {
				return ScreenMsg(LoginScreen)
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
						return ScreenMsg(LoginScreen)
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

func (m *SignupModel) View() string {
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

func (m *SignupModel) handleFocus() {
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}
