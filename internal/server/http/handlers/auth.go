package controllers

import (
	"context"
	"errors"

	"github.com/K44Z/kzchat/internal/server/http"
	"github.com/K44Z/kzchat/internal/server/schemas"
	"github.com/K44Z/kzchat/internal/server/services"

	repository "github.com/K44Z/kzchat/internal/server/database/generated"

	"github.com/K44Z/kzchat/internal/helpers"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := c.Locals("validatedBody").(schemas.Auth)
		ctx := context.Background()
		_, err := s.UserService.GetUserByUsername(ctx, body.Username)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				hashedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
				if err != nil {
					log.Printf("Error hashing password: %v\n", err)
					return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
				}
				user := repository.CreateUserParams{
					Username: body.Username,
					Password: string(hashedPass),
				}
				err = s.UserService.CreateUser(ctx, user.Username, user.Password)
				if err != nil {
					log.Printf("Error creating user: %v\n", err)
					return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
				}
				return http.Success(c, nil)
			}
			log.Printf("DB error: %v\n", err)
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		return http.Error(c, fiber.ErrBadRequest.Code, "Username already taken")
	}
}

func LoginHanlder(s *services.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := c.Locals("validatedBody").(schemas.Auth)
		ctx := context.Background()
		user, err := s.UserService.GetUserWithPassword(ctx, body.Username)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println(err)
			return http.Error(c, fiber.ErrNotFound.Code, fiber.ErrNotFound.Error())
		} else if err != nil {
			log.Println(err)
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			log.Println(err)
			return http.Error(c, fiber.ErrUnauthorized.Code, "Invalid credentials")
		}
		token, err := helpers.GenerateJWTtoken(schemas.User{
			Username: user.Username,
			ID:       user.ID,
		})
		if err != nil {
			log.Println(err)
			return http.Error(c, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Error())
		}
		return http.Success(c, map[string]any{
			"token": token,
		})
	}
}
