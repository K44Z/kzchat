package controllers

import (
	"context"

	"github.com/K44Z/kzchat/internal/server/http"
	"github.com/K44Z/kzchat/internal/server/services"

	"github.com/gofiber/fiber/v2"
)

func SearchForUserByUsernameHandler(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		if username == "" {
			return http.Error(c, fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Error())
		}
		_, err := s.UserService.GetUserByUsername(c.Context(), username)
		if err != nil {
			return http.Error(c, fiber.ErrNotFound.Code, fiber.ErrNotFound.Error())
		}
		return http.Success(c, nil)
	}
}

func GetAllUsersHandler(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		users, err := s.UserService.GetAllUsersService(ctx)
		if err != nil {
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		return http.Success(c, map[string]interface{}{
			"users": users,
		})
	}
}
