package middleware

import (
	"log"
	"strings"

	"github.com/K44Z/kzchat/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func Jwt(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		log.Println("Missing authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing or invalid Authorization header",
		})
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("invalid Authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := helpers.Authenticate(tokenString)
	if err != nil {
		log.Println("Invalid or expired token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}
	c.Locals("user", user)
	return c.Next()
}
