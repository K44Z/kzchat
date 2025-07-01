package routes

import (
	h "kzchat/server/handlers"
	"kzchat/server/middleware"
	"kzchat/server/schemas"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(router fiber.Router) {
	router.Post("/register", middleware.ValidateBody[schemas.Auth](), h.Register)
	router.Post("/login", middleware.ValidateBody[schemas.Auth](), h.Login)
}
