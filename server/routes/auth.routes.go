package routes

import (
	"kzchat/server/controllers"
	"kzchat/server/schemas"
	"kzchat/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(router fiber.Router) {
	router.Post("/register", middleware.ValidateBody[schemas.Auth](), controllers.Register)
	router.Post("/login", middleware.ValidateBody[schemas.Auth](),controllers.Login)
}
