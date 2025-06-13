package controllers

import (
	"context"
	"errors"
	"fmt"
	"kzchat/server/database"
	repository "kzchat/server/database/generated"
	"kzchat/server/schemas"

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

func CreateDmMessage(c *fiber.Ctx) error {
	ctx := context.Background()
	recUsername := c.Params("username")
	body := c.Locals("validatedBody").(schemas.CreateMessageSchema)
	user := c.Locals("user").(repository.User)
	if body.ChatId == 0 {
		rec, err := database.Queries.GetUserByUsername(ctx, recUsername)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Recipient not found",
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal server error",
				})
			}
		}
		name := fmt.Sprintf("%s - %s ", user.Username, recUsername)
		fmt.Println(name)
		chatParams := repository.CreateChatParams{
			Type: body.Type,
			Name: name,
		}
		chat, err := database.Queries.CreateChat(ctx, chatParams)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}
		_, err = database.Queries.CreateChatMembers(ctx, repository.CreateChatMembersParams{
			ChatID: chat.ID,
			UserID: user.ID,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}
		_, err = database.Queries.CreateChatMembers(ctx, repository.CreateChatMembersParams{
			ChatID: body.ChatId,
			UserID: rec.ID,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}
	}
	timestamp := pgtype.Timestamp{
		Time:  body.Time,
		Valid: true,
	}

	params := repository.StoreChatMessageParams{
		SenderID: body.SenderId,
		Content:  body.Content,
		ChatID:   body.ChatId,
		Time:     timestamp,
		Type:     body.Type,
	}
	database.Queries.StoreChatMessage(ctx, params)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Message sent successfully",
	})
}

// func CreateChanMessage(c *fiber.Ctx) error {
// 	ctx := context.Background()
// 	chanName := c.Params("name")
// 	body := c.Locals("validatedBody").(schemas.CreateMessageSchema)	
// 	return c.SendStatus(fiber.StatusCreated)
// }
 

// func CreateChat() error {
	
// }