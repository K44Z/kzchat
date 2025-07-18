package routes

import (
	"github.com/K44Z/kzchat/internal/server/middleware"
	"github.com/K44Z/kzchat/internal/server/schemas"

	h "github.com/K44Z/kzchat/internal/server/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(router fiber.Router) {
	router.Post("/register", middleware.ValidateBody[schemas.Auth](), h.Register)
	router.Post("/login", middleware.ValidateBody[schemas.Auth](), h.Login)
}
