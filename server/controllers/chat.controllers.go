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

func GetDmsByrecipientUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	ctx := context.Background()
	current := c.Locals("user").(repository.User)
	rec, err := database.Queries.GetUserByUsername(ctx, username)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprint("User ", username, " does not exist"),
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	fmt.Println("current :", current)
	params := repository.GetDmChatMessagesByParticipantsParams{
		UserID:   current.ID,
		UserID_2: rec.ID,
	}
	messages, err := database.Queries.GetDmChatMessagesByParticipants(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	fmt.Println("messages: ", messages)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": messages,
	})
}

func CreateDmMessage(m schemas.Message) error {
	ctx := context.Background()
	var (
		chat      *repository.Chat
		err       error
		usernames = []string{
			m.SenderUsername, m.ReceiverUsername,
		}
		users []repository.User
	)

	for _, u := range usernames {
		user, err := database.Queries.GetUserByUsername(ctx, u)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("user does not exist")
			} else {
				return fmt.Errorf("internal server error : %s", err)
			}
		}
		users = append(users, user)
	}
	if m.Chat.ID == 0 {
		chat, err = services.CreateDMChatFromMessage(m, users)
		if err != nil {
			return fmt.Errorf("error creating chat :%s", err)
		}
	}

	timestamp := pgtype.Timestamp{
		Time:  m.Time,
		Valid: true,
	}
	params := repository.StoreChatMessageParams{
		SenderID: users[0].ID,
		Content:  m.Content,
		ChatID:   chat.ID,
		Time:     timestamp,
		Type:     "dm",
	}
	_, err = database.Queries.StoreChatMessage(ctx, params)
	if err != nil {
		return fmt.Errorf("error storing message :%s", err)
	}
	return nil
}

func GetChatByParticipants(c *fiber.Ctx) error {
	ctx := context.Background()
	var users []schemas.User
	fmt.Println(c.Locals("validatedBody"))
	body, ok := c.Locals("validatedBody").(schemas.GetChatIdByParticipants)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to validated body",
		})
	}
	for _, u := range body.Members {
		user, err := database.Queries.GetUserByUsername(ctx, u)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("user does not exist")
			} else {
				return fmt.Errorf("internal server error : %s", err)
			}
		}
		users = append(users, schemas.User{
			Username: user.Username,
			ID:       user.ID,
		})
	}
	chatId, err := database.Queries.FindChatByParticipants(ctx, repository.FindChatByParticipantsParams{
		Column1: []string{users[0].Username, string(users[1].Username)},
		UserID:  2, // this is the expected count of users not userid
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Chat not found",
			})
		}
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	fmt.Println("the response is :", chatId, users)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"chatId": chatId,
		"users":  users,
	})
}
