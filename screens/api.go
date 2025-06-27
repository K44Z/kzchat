package screens

import (
	"encoding/json"
	"io"
	a "kzchat/auth"
	repository "kzchat/server/database/generated"
	"kzchat/server/schemas"
	"net/http"

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

type Screen int
type ScreenMsg Screen

const (
	SignupScreen Screen = iota
	LoginScreen
	ChatScreen
)

var Messages = make(chan tea.Msg)

func (m ChatModel) ConnectToWs() tea.Cmd {
	return func() tea.Msg {
		c, _, err := websocket.DefaultDialer.Dial(WS_URL+"/ws", nil)
		if err != nil {
			return ErrMsg(err)
		}
		return WsConnectedMsg{c}
	}
}

func (m ChatModel) ReadLoop(conn *websocket.Conn) {
	for {
		var msg schemas.Message
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		Messages <- WsMsg{Content: msg.Content, Time: msg.Time, SenderUsername: msg.SenderUsername}
	}
}

func (c ChatModel) FetchMessages(recipientId string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest("GET", API_URL+"/message/recId/"+recipientId, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+a.Config.Token)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var messages []repository.Message
		if err := json.Unmarshal(body, &messages); err != nil {
			return err
		}
		// return ChatFetchedMsg{Messages: messages}
		return nil
	}
}
