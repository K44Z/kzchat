package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	AuthRouter(app.Group("/auth"))
	UserRouter(app.Group("/users"))
	MessagesRouter(app.Group("/messages"))
	SocketRouter(app.Group("/ws"))
}
