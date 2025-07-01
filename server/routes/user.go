package routes

import (
	h "kzchat/server/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router) {
	router.Get("/:username", h.SearchForUserByUsername)
}
