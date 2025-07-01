package controllers

import (
	"kzchat/server/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SearchForUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	_, err := services.GetUserByUsername(username)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user not found",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
