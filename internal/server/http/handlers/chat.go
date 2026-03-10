package controllers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/K44Z/kzchat/internal/server/http"

	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/K44Z/kzchat/internal/server/services"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func GetDmsByrecipientUsernameHandler(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		if username == "" {
			log.Println("Username is required")
			return http.Error(c, fiber.ErrBadRequest.Code, "Username is required")
		}
		currentUser := c.Locals("user").(*schemas.User)
		rec, err := s.UserService.GetUserByUsername(c.Context(), username)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println(err)
			return http.Error(c, fiber.ErrNotFound.Code, fmt.Sprint("User ", username, " does not exist"))
		} else if err != nil {
			log.Println(err)
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		messages, err := s.ChatService.GetMessagesByParticipants(c.Context(), *currentUser, *rec)
		if err != nil {
			log.Println(err)
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		return http.Success(c, map[string]interface{}{
			"messages": messages,
		})
	}
}

func GetChatByParticipantsHandler(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []schemas.User
		body := c.Locals("validatedBody").(schemas.GetChatIdByParticipants)
		for _, u := range body.Members {
			user, err := s.UserService.GetUserByUsername(c.Context(), u)
			if err != nil {
				log.Println(err)
				if errors.Is(err, pgx.ErrNoRows) {
					return http.Error(c, fiber.ErrNotFound.Code, fiber.ErrNotFound.Error())
				} else {
					return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
				}
			}
			users = append(users, schemas.User{
				Username: user.Username,
				ID:       user.ID,
			})
		}
		chatID, name, _ := s.ChatService.GetChatIdByParticipants(c.Context(), []int32{users[0].ID, users[1].ID})
		if chatID == nil {
			chat, err := s.ChatService.CreateChat(c.Context(), users, "dm")
			if err != nil {
				log.Println(err)
				return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
			}
			chatID = &chat.ID
			name = chat.Name
		}
		attachements, err := s.ChatService.GetAttachmentsByChatService(c.Context(), *chatID)
		if err != nil {
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		return http.Success(c, fiber.Map{
			"chatId":      *chatID,
			"name":        name,
			"attachments": attachements,
			"users":       users,
		})
	}
}

func CreateChatFromMessageHandler(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []schemas.User
		body := c.Locals("validatedBody").(schemas.CreateChatByMessage)
		for _, u := range body.Members {
			user, err := s.UserService.GetUserByUsername(c.Context(), u)
			if err != nil {
				log.Print(err)
				if errors.Is(err, pgx.ErrNoRows) {
					return http.Error(c, fiber.ErrNotFound.Code, fiber.ErrNotFound.Error())
				} else {
					return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
				}
			}
			users = append(users, *user)
		}
		chat, err := s.ChatService.CreateChatFromMessage(c.Context(), body.Message, users)
		if err != nil {
			log.Print(err)
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		return http.Created(c, map[string]any{
			"chat": chat,
		})
	}
}

func UploadHandler(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		chatID := c.Params("id")
		if chatID == "" {
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		f, err := c.FormFile("file")
		if err != nil {
			log.Println(err)
			return http.Error(c, 500, fiber.ErrInternalServerError.Error())
		}
		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create upload directory: %w", err)
		}
		path := filepath.Join(uploadDir, f.Filename)
		err = c.SaveFile(f, path)
		if err != nil {
			log.Println(err)
			return http.Error(c, 500, fiber.ErrInternalServerError.Error())
		}
		id, err := strconv.Atoi(chatID)
		if err != nil {
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		err = s.ChatService.SaveAttachment(c.Context(), schemas.SaveAttachementParams{
			File:   f,
			ChatID: int32(id),
			Path:   path,
		})
		if err != nil {
			return http.Error(c, 500, fiber.ErrInternalServerError.Error())
		}
		return http.Success(c, map[string]any{
			"filename": f.Filename,
		})
	}
}

func GetAttachmentsByChatHandler(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		chatID := c.Params("id")
		if chatID == "" {
			return http.Error(c, fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Error())
		}
		id, err := strconv.Atoi(chatID)
		if err != nil {
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		attachments, err := s.ChatService.GetAttachmentsByChatService(c.Context(), int32(id))
		if err != nil {
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		return http.Success(c, map[string]any{
			"attachments": attachments,
		})
	}
}

func DownloadFileHandler(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := c.Locals("validatedBody").(schemas.DownloadFile)
		return c.Status(200).Download(req.Path)
	}
}
