package controllers

import (
	"fmt"
	"kzchat/auth"
	"kzchat/helpers"
	"kzchat/server/models"
	"kzchat/server/schemas"
	"kzchat/server/services"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(schemas.Auth)
	exists,err := services.CheckExistingUser(body.Username)
	if err != nil || exists {
		helpers.Logger.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already exists",
		})
	}	
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.Logger.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}
	user := &models.User{
		Username: body.Username,
		Password: string(hashedPass),
	}
	err = services.CreateUser(user)
	if err != nil {
		helpers.Logger.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error while creating user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
	})
}

func Login(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(schemas.Auth)
	user, err := services.GetUserByUsername(body.Username)
	if err != nil {
		helpers.Logger.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User does not exist" + err.Error(),
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login Sucessfull",
		"token":   token,
	})
}
