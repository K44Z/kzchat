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
	router.Post("/chat", middleware.Jwt, middleware.ValidateBody[schemas.GetChatIdByParticipants](), h.GetChatByParticipantsHandler(s))
	router.Post("/createChat", middleware.Jwt, middleware.ValidateBody[schemas.CreateChatByMessage](), h.CreateChatFromMessageHandler(s))
	router.Post("/upload", middleware.Jwt, h.UploadHandler(s))
}
