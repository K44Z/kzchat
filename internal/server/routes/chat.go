package routes

import (
	h "github.com/K44Z/kzchat/internal/server/handlers"
	"github.com/K44Z/kzchat/internal/server/middleware"
	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/gofiber/fiber/v2"
)

func MessagesRouter(router fiber.Router) {
	router.Get("/recipient/:username", middleware.Jwt, h.GetDmsByrecipientUsername)
	router.Post("/chat", middleware.ValidateBody[schemas.GetChatIdByParticipants](), middleware.Jwt, h.GetChatByParticipants)
	router.Post("/createChat", middleware.ValidateBody[schemas.CreateChatByMessage](), middleware.Jwt, h.CreateChatFromMessage)
}
