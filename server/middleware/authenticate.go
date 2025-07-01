package middleware

import (
	"kzchat/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing or invalid Authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := helpers.Authenticate(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}
	c.Locals("user", user)
	return c.Next()
}
