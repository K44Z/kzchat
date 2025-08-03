package routes

import (
	"github.com/K44Z/kzchat/internal/server/http/middleware"
	"github.com/K44Z/kzchat/internal/server/schemas"
	"github.com/K44Z/kzchat/internal/server/services"

	h "github.com/K44Z/kzchat/internal/server/http/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(router fiber.Router, s *services.Services) {
	router.Post("/register", middleware.ValidateBody[schemas.Auth](), h.RegisterHandler(s))
	router.Post("/login", middleware.ValidateBody[schemas.Auth](), h.LoginHanlder(s))
}
