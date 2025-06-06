package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	AuthRouter(app.Group("/auth"))	
	MessagesRouter(app.Group("/messages"))
	SocketRouter(app.Group("/ws"))
	// app.Get("/swagger/*", swagger.HandlerDefault)
	// app.Use("/profile", ProfileRouter)
}
