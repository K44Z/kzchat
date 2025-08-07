package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, data map[string]interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Status: "success",
		Data:   data,
	})
}

func Created(c *fiber.Ctx, data map[string]interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Status: "success",
		Data:   data,
	})
}

func Error(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(APIResponse{
		Status:  strconv.Itoa(status),
		Message: msg,
	})
}
