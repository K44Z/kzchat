package routes

import (
	h "github.com/K44Z/kzchat/internal/server/http/handlers"
	"github.com/K44Z/kzchat/internal/server/http/middleware"
	"github.com/K44Z/kzchat/internal/server/schemas"
	"github.com/K44Z/kzchat/internal/server/services"

	"github.com/gofiber/fiber/v2"
)

func MessagesRouter(router fiber.Router, s *services.Services) {
	router.Get("/recipient/:username", middleware.Jwt, h.GetDmsByrecipientUsernameHandler(s))
	router.Post("/chat", middleware.ValidateBody[schemas.GetChatIdByParticipants](), middleware.Jwt, h.GetChatByParticipantsHandler(s))
	router.Post("/createChat", middleware.ValidateBody[schemas.CreateChatByMessage](), middleware.Jwt, h.CreateChatFromMessageHandler(s))
}
