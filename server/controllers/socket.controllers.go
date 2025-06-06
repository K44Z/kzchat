package controllers

import (
	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)

func Broadcast(c *websocket.Conn) {
	clients[c] = true
	defer func() {
		delete(clients, c)
		c.Close()
	}()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
