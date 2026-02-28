package messages

import (
	"github.com/K44Z/kzchat/internal/server/schemas"
	"github.com/gorilla/websocket"
)

type (
	Screen    int
	ScreenMsg Screen
	ErrMsg    error
	WsMsg     schemas.Message
)

type WsConnectedMsg struct {
	Conn *websocket.Conn
}

type ChatFetchedMsg struct {
	Messages []schemas.Message
}
