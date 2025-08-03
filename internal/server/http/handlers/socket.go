package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/K44Z/kzchat/internal/server/schemas"
	"github.com/K44Z/kzchat/internal/server/services"

	"github.com/K44Z/kzchat/internal/helpers"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)

func Broadcast(s *services.Services) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
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
			fmt.Println("message received:", string(msg))

			var m schemas.Message
			err = json.Unmarshal(msg, &m)
			if err != nil {
				fmt.Printf("error unmarshal json: %s", err)
				continue
			}

			err = helpers.ValidateStruct(m) // todo send errors
			if err != nil {
				log.Error("Error validating struct", err)
				continue
			}

			if err = s.ChatService.CreateDM(context.Background(), m); err != nil {
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
}
