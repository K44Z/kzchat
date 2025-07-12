package controllers

import (
	"encoding/json"
	"fmt"
	"kzchat/helpers"
	"kzchat/server/schemas"

	"github.com/gofiber/fiber/v2/log"
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
		fmt.Println("message received :", string(msg))
		var m schemas.Message
		err = json.Unmarshal(msg, &m)
		if err != nil {
			fmt.Printf("error unmarshal json: %s", err)
		}
		err = helpers.ValidateStruct(m) // todo send errors
		if err != nil {
			log.Error("Error vaidating strcut", err)
			continue
		}
		if err = CreateDmMessage(m); err != nil {
			log.Error("error creating message", err)
			continue
		}
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
