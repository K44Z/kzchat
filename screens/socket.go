package screens

import (
	"encoding/json"
	"io"
	"kzchat/api"
	a "kzchat/api"
	"kzchat/server/schemas"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

const (
	API_URL = "http://localhost:4000"
	WS_URL  = "ws://localhost:4000"
)

type WsMsg schemas.Message
type ErrMsg error
type WsConnectedMsg struct {
	Conn *websocket.Conn
}

type ChatFetchedMsg struct {
	Messages []schemas.Message
}
type Screen int
type ScreenMsg Screen

const (
	SignupScreen Screen = iota
	LoginScreen
	ChatScreen
)

var Messages = make(chan tea.Msg)

func (m *ChatModel) ConnectToWs() tea.Cmd {
	return func() tea.Msg {
		c, _, err := websocket.DefaultDialer.Dial(WS_URL+"/ws", nil)
		if err != nil {
			return ErrMsg(err)
		}
		return WsConnectedMsg{c}
	}
}

func (m *ChatModel) ReadLoop(conn *websocket.Conn) {
	for {
		var msg schemas.Message
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		Messages <- WsMsg{Content: msg.Content, Time: msg.Time, SenderUsername: msg.SenderUsername}
	}
}

func (m *ChatModel) FetchMessages() tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest("GET", API_URL+"/messages/recipient/"+string(m.Recipient.Username), nil)
		if err != nil {
			return ErrMsg(err)
		}
		req.Header.Add("Authorization", "Bearer "+a.Config.Token)
		resp, err := client.Do(req)
		if err != nil {
			return ErrMsg(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ErrMsg(err)
		}
		var messages []schemas.Message
		if err := json.Unmarshal(body, &messages); err != nil {
			return ErrMsg(err)
		}
		m.Messages = messages
		return ChatFetchedMsg{Messages: messages}
	}
}

func (m *ChatModel) SendMessage() {
	inputValue := strings.TrimSpace(m.Input.Value())
	if m.Recipient.ID == 0 || m.Recipient.Username == "" {
		m.Err = "Please select a recipient before sending a message"
		return
	}
	if inputValue != "" {
		message := schemas.Message{
			Content:          inputValue,
			SenderUsername:   api.Config.Username,
			Time:             time.Now(),
			Chat:             m.Chat,
			ReceiverUsername: m.Recipient.Username,
		}
		if m.Ws == nil {
			if m.Err != "" {
				m.Err += "; ws is not connected"
			} else {
				m.Err = "ws is not connected"
			}
			return
		}
		err := m.Ws.WriteJSON(message)
		if err != nil {
			m.Err = err.Error()
		}
		m.Messages = append(m.Messages, message)
		m.Viewport.SetContent(m.renderMessages())
		m.Viewport.GotoBottom()
		m.Input.Reset()
	}
}
