package routes

import (
	h "kzchat/server/handlers"
	"kzchat/server/middleware"
	"kzchat/server/schemas"

	"github.com/gofiber/fiber/v2"
)

func MessagesRouter(router fiber.Router) {
	router.Get("/recipient/:username", middleware.JWTMiddleware, h.GetDmsByrecipientUsername)
	router.Post("/chat", middleware.ValidateBody[schemas.GetChatIdByParticipants](), h.GetChatByParticipants)
	// router.Post("/channel/:name", controllers.CreateChanMessage)
}
