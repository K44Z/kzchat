package main

import (
	"encoding/json"
	"io"
	repository "kzchat/server/database/generated"
	"kzchat/server/schemas"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

var API_URL = "http://localhost:4000"
var WS_URL = "ws://localhost:4000"

type wsMsg schemas.Message
type errMsg error
type wsConnectedMsg struct {
	conn *websocket.Conn
}
var messages = make(chan tea.Msg)

func (m ChatModel) connectToWs() tea.Cmd {
	return func() tea.Msg {
		c, _, err := websocket.DefaultDialer.Dial(WS_URL+"/ws", nil)
		if err != nil {
			return errMsg(err)
		}
		return wsConnectedMsg{c}
	}
}

func (m ChatModel) readLoop(conn *websocket.Conn) {
	for {
		var msg schemas.Message
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		messages <- wsMsg{Content: msg.Content, Time: msg.Time, SenderUsername: msg.SenderUsername}
	}
}

func (c ChatModel) fetchMessages(recipientId string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest("GET", API_URL+"/message/recId/"+recipientId, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+config.Token)
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
