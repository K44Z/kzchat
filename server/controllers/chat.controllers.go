package controllers

import (
	"context"
	"database/sql"
	"errors"
	"kzchat/server/database"
	repository "kzchat/server/database/generated"

	"github.com/gofiber/fiber/v2"
)

func GetMessages(c *fiber.Ctx) error {
	return nil
}

func GetDmByrecipientUsername(c *fiber.Ctx) error {
	username := c.Params("recUsername")
	ctx := context.Background()
	user, err := database.Queries.GetUserByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already exists",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	params := repository.GetDmChatMessagesByParticipantsParams{
		UserID:   int32(c.Locals("id").(float64)),
		UserID_2: user.ID,
	}
	messages, err := database.Queries.GetDmChatMessagesByParticipants(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": messages,
	})
}
