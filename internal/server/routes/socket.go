package routes

import (
	h "github.com/K44Z/kzchat/internal/server/handlers"
	"github.com/K44Z/kzchat/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SocketRouter(router fiber.Router) {
	router.Get("/", middleware.Jwt, websocket.New(h.Broadcast))
}
