package middleware

import (
	"log"

	"github.com/K44Z/kzchat/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		if err := c.BodyParser(&body); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Invalid request body",
			})
		}

		if err := helpers.ValidateStruct(&body); err != nil {
			log.Print(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Invalid request body",
			})
		}
		log.Println("BODY VALIDATED")
		c.Locals("validatedBody", body)
		return c.Next()
	}
}
