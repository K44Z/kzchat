package controllers

import (
	"context"
	"errors"
	"fmt"
	"kzchat/helpers"
	"kzchat/server/database"
	repository "kzchat/server/database/generated"
	"kzchat/server/schemas"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(schemas.Auth)
	ctx := context.Background()
	_, err := database.Queries.GetUserByUsername(ctx, body.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			hashedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to hash password",
				})
			}
			user := repository.CreateUserParams{
				Username: body.Username,
				Password: string(hashedPass),
			}
			str, err := database.Queries.CreateUser(ctx, user)
			if err != nil {
				log.Println("2", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal server error",
				})
			}
			log.Println(str)
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"message": "User created",
			})
		}
		log.Printf("DB error: %T\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "Username already taken",
	})
}

func Login(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(schemas.Auth)
	ctx := context.Background()
	user, err := database.Queries.GetUserByUsername(ctx, body.Username)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already exists",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	fmt.Println(user.Password)
	fmt.Println(body.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{

			"message": "Invalid credentials",
		})
	}
	token, err := helpers.GenerateJWTtoken(user)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login pgx",
		"token":   token,
	})
}
