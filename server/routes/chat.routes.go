package routes

import (
	"kzchat/server/controllers"

	"github.com/gofiber/fiber/v2"
)

func MessagesRouter(router fiber.Router) {
	router.Get("/recipient/:username", controllers.GetDmByrecipientUsername)
	// router.Post("/channel/:name", controllers.CreateChanMessage)
}
