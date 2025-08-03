package routes

import (
	"github.com/K44Z/kzchat/internal/server/services"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, s *services.Services) {
	AuthRouter(app.Group("/auth"), s)
	UserRouter(app.Group("/users"), s)
	MessagesRouter(app.Group("/messages"), s)
	SocketRouter(app.Group("/ws"), s)
}
