package middleware

import (
	"kzchat/helpers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		if err := c.BodyParser(&body); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"message": "Invalid request body",
			})
		}

		if err := helpers.ValidateStruct(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"message": "Invalid credentials please ensure password is 6 char long",
			})
		}
		log.Println("BODY VALIDATED")
		c.Locals("validatedBody", body)
		return c.Next()
	}
}