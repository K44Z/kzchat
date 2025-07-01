package routes

import (
	"kzchat/server/controllers"
	"kzchat/server/middleware"
	"kzchat/server/schemas"

	"github.com/gofiber/fiber/v2"
)

func MessagesRouter(router fiber.Router) {
	router.Get("/recipient/:username", middleware.JWTMiddleware, controllers.GetDmsByrecipientUsername) 
	router.Post("/chatId", middleware.ValidateBody[schemas.GetChatIdByParticipants](), controllers.GetChatByParticipants)
	// router.Post("/channel/:name", controllers.CreateChanMessage)
}
