package routes

import (
	c "github.com/K44Z/kzchat/internal/server/http/handlers"
	"github.com/K44Z/kzchat/internal/server/http/middleware"
	"github.com/K44Z/kzchat/internal/server/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SocketRouter(router fiber.Router, s *services.Services) {
	router.Get("/", middleware.Jwt, websocket.New(c.Broadcast(s)))
}
