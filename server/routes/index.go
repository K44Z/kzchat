package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	AuthRouter(app.Group("/auth"))	
	// app.Get("/swagger/*", swagger.HandlerDefault)
	MessagesRouter(app.Group("/messages"))
	// app.Use("/profile", ProfileRouter)
}
