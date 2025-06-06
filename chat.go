package main

import (
	"fmt"
	repository "kzchat/server/database/generated"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type ChatModel struct {
	messages    []repository.Message
	message     string
	commandMode bool
	command     string
	input       textinput.Model
	username    string
	err         string
	channels    []string
	width       int
	ws          *websocket.Conn
	height      int
	viewport    viewport.Model
	content     string
	ready       bool
}

func NewChatModel(username string) ChatModel {
	input := textinput.New()
	input.Prompt = ""
	input.Focus()
	input.CharLimit = 500
	input.Width = 65
	return ChatModel{
		messages: []repository.Message{},
		username: username,
		input:    input,
		channels: []string{"general", "random", "dev", "help"},
	}
}

func (m ChatModel) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m ChatModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m ChatModel) Init() tea.Cmd {
	return receive(m.ws)
}

func (m ChatModel) Update(msg tea.Msg) (ChatModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.message = m.input.Value()
			// m.sendMessage()
			m.input.Reset()
		}
	}
	if !m.commandMode {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m ChatModel) View() string {

	leftWidth := int(float64(m.width) * 0.14)
	rightWidth := int(float64(m.width) * 0.14)
	chatWidth := m.width - leftWidth - rightWidth - 7
	contentHeight := m.height - 7

	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(contentHeight).
		Padding(1, 1).
		Align(lipgloss.Left)

	chatStyle := lipgloss.NewStyle().
		Width(chatWidth).
		Height(contentHeight).
		Padding(1, 1).
		Align(lipgloss.Left).
		Border(lipgloss.NormalBorder(), false, true, false, true).
		BorderForeground(lipgloss.Color("240"))

	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(contentHeight).
		Padding(1, 1).
		Align(lipgloss.Left)

	widthString := fmt.Sprint(m.width)
	heightString := fmt.Sprint(m.height)

	leftSection := leftStyle.Render(m.renderLeftSidebar())
	chatSection := chatStyle.Render(m.renderMessages())
	rightSection := rightStyle.Render(m.renderRightSidebar() + widthString + " " + heightString)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, chatSection, rightSection)

	inputStyle := lipgloss.NewStyle().
		Width(m.width-4).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("240"))

	inputSection := inputStyle.Render(m.input.View())

	var content strings.Builder
	content.WriteString(mainContent)
	content.WriteString(inputSection)

	if m.err != "" {
		errorStyle := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("9")).
			Bold(true)

		errorMsg := errorStyle.Render("Error: " + m.err)
		content.WriteString("\n\n" + errorMsg)
	}
	return content.String()
}
