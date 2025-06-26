package routes

import (
	"kzchat/server/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router) {
	router.Get("/:username", controllers.SearchForUserByUsername)
}
