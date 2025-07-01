package routes

import (
	h "kzchat/server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SocketRouter(router fiber.Router) {
	router.Get("/", websocket.New(h.Broadcast))
}
