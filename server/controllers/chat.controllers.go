package controllers

import (
	"context"
	"errors"
	"fmt"
	"kzchat/server/database"
	repository "kzchat/server/database/generated"
	"kzchat/server/schemas"
	"kzchat/server/services"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetMessages(c *fiber.Ctx) error {
	return nil
}

func GetDmByrecipientUsername(c *fiber.Ctx) error {
	username := c.Params("recUsername")
	ctx := context.Background()
	user, err := database.Queries.GetUserByUsername(ctx, username)
	if errors.Is(err, pgx.ErrNoRows) {
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

func CreateDmMessage(m schemas.Message) error {
	ctx := context.Background()
	var chat *repository.Chat
	rec, err := database.Queries.GetUserByUsername(ctx, m.ReceiverUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user does not exist")
		} else {
			return fmt.Errorf("internal server error : %s", err)
		}
	}
	user, err := database.Queries.GetUserByUsername(ctx, m.SenderUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user does not exist")
		} else {
			return fmt.Errorf("internal server error : %s", err)
		}
	}

	if m.Chat.ID == 0 {
		users := []repository.User{
			user, rec,
		}
		fmt.Println("the users are",users)
		chat, err = services.CreateDMChatFromMessage(m, users)
		if err != nil {
			return err
		}
	}

	timestamp := pgtype.Timestamp{
		Time:  m.Time,
		Valid: true,
	}
	params := repository.StoreChatMessageParams{
		SenderID: user.ID,
		Content:  m.Content,
		ChatID:   chat.ID,
		Time:     timestamp,
		Type:     "dm",
	}
	database.Queries.StoreChatMessage(ctx, params)
	return nil
}
