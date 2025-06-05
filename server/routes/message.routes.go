package routes

import (
	"kzchat/server/controllers"

	"github.com/gofiber/fiber/v2"
)

func MessagesRouter(router fiber.Router) {
	router.Get("/message/recipientId/:id", controllers.GetMessages)
}
