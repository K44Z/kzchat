package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"

	repository "github.com/K44Z/kzchat/internal/server/database/generated"

	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/K44Z/kzchat/internal/server/services"

	"github.com/K44Z/kzchat/internal/server/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

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
	users := []repository.User{
		current, rec,
	}
	params := repository.GetDmChatMessagesByParticipantsParams{
		UserID:   current.ID,
		UserID_2: rec.ID,
	}
	m, err := database.Queries.GetDmChatMessagesByParticipants(ctx, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	messages, err := services.MapMessagesToClient(m, users)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(messages)
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
	fmt.Println("the chat id is :", m.Chat.ID)
	if m.Chat.ID == 0 {
		chat, err = services.CreateDMChatFromMessage(m, users)
		if err != nil {
			return fmt.Errorf("error creating chat :%s", err)
		}
	} else {
		chat = &repository.Chat{ID: m.Chat.ID}
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
	body := c.Locals("validatedBody").(schemas.GetChatIdByParticipants)
	for _, u := range body.Members {
		user, err := database.Queries.GetUserByUsername(ctx, u)
		if err != nil {
			log.Println(err)
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
		Column1: []int32{users[0].ID, users[1].ID},
		Column2: 2,
	})
	if chatId == 0 {
		log.Printf("Chat not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Chat not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	fmt.Println("the chat id is :", chatId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"chatId": chatId,
		"users":  users,
	})
}

func CreateChatFromMessage(c *fiber.Ctx) error {
	ctx := context.Background()
	var users []repository.User
	body := c.Locals("validatedBody").(schemas.CreateChatByMessage)
	for _, u := range body.Members {
		user, err := database.Queries.GetUserByUsername(ctx, u)
		if err != nil {
			log.Print(err)
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("user does not exist")
			} else {
				return fmt.Errorf("internal server error : %s", err)
			}
		}
		users = append(users, user)
	}
	chat, err := services.CreateDMChatFromMessage(body.Message, users)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"chat": chat,
	})
}
