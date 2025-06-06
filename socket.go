package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)


type wsMsg string 
type errMsg error

func receive(ws *websocket.Conn) func() tea.Msg {
	return func () tea.Msg {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return errMsg(err)
		}
		return wsMsg(msg)
	}
}
