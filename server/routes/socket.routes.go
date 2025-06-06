package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"kzchat/server/controllers"
)

func SocketRouter(router fiber.Router) {
	router.Get("/", websocket.New(controllers.Broadcast))
}
