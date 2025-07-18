package routes

import (
	h "github.com/K44Z/kzchat/internal/server/handlers"
	"github.com/K44Z/kzchat/internal/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router) {
	router.Get("/:username", middleware.Jwt, h.SearchForUserByUsername)
}
