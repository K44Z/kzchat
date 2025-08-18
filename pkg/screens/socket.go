package screens

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/K44Z/kzchat/internal/api"

	"github.com/K44Z/kzchat/internal/server/schemas"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

const (
	SignupScreen api.Screen = iota
	LoginScreen
	ChatScreen
)

var Messages = make(chan tea.Msg)

type FetchMessagesResponse struct {
	Status string              `json:"status"`
	Data   MessagesReponseData `json:"data"`
}

type MessagesReponseData struct {
	Messages []schemas.Message `json:"messages"`
}

func (m *ChatModel) ConnectToWs() tea.Cmd {
	return func() tea.Msg {
		header := http.Header{}
		header.Add("Authorization", "Bearer "+api.Config.Token)
		c, _, err := websocket.DefaultDialer.Dial(api.WS_URL+"/ws", header)
		if err != nil {
			return api.ErrMsg(err)
		}
		return api.WsConnectedMsg{c}
	}
}

func (m *ChatModel) ReadLoop(conn *websocket.Conn) {
	for {
		var msg schemas.Message
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		Messages <- api.WsMsg{Content: msg.Content, Time: msg.Time, Sender: msg.Sender}
	}
}

func (m *ChatModel) FetchMessages() tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest("GET", api.BASE_URL+"/messages/recipient/"+string(m.Recipient.Username), nil)
		if err != nil {
			return api.ErrMsg(err)
		}
		req.Header.Add("Authorization", "Bearer "+api.Config.Token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return api.ErrMsg(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return api.ErrMsg(err)
		}
		var apiResp FetchMessagesResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return api.ErrMsg(err)
		}
		return api.ChatFetchedMsg{Messages: apiResp.Data.Messages}
	}
}

func (m *ChatModel) SendMessage(recipient schemas.User) {
	inputValue := strings.TrimSpace(m.Textarea.Value())
	if m.Recipient.ID == 0 || m.Recipient.Username == "" {
		m.Err = "Please select a recipient before sending a message"
		return
	}
	if inputValue != "" {
		message := schemas.Message{
			Content: inputValue,
			Sender: schemas.User{
				Username: api.Config.Username,
			},
			Time:     time.Now(),
			Chat:     m.Chat,
			Receiver: m.Recipient,
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
		if m.Recipient == recipient {
			m.Messages = append(m.Messages, message)
		}
		m.Viewport.SetContent(m.RenderMessages())
		m.Viewport.GotoBottom()
		m.Textarea.Reset()
	}
}
