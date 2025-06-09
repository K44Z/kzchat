package controllers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	authentication "kzchat/auth"
	"kzchat/helpers"
	"kzchat/server/database"
	repository "kzchat/server/database/generated"
	"kzchat/server/schemas"

	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(schemas.Auth)
	ctx := context.Background()
	_, err := database.Queries.GetUserByUsername(ctx, body.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already exists",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.Logger.Println(err)
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	log.Println(str)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
	})
}

func Login(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(schemas.Auth)
	ctx := context.Background()
	user, err := database.Queries.GetUserByUsername(ctx, body.Username)
	if errors.Is(err, sql.ErrNoRows) {
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
		helpers.Logger.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}
	token, err := authentication.GenerateJWTtoken(user)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login Sucessfull",
		"token":   token,
	})
}
