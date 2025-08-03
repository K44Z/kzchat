package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/K44Z/kzchat/configs"
	"github.com/K44Z/kzchat/internal/server/http/routes"
	"github.com/K44Z/kzchat/internal/server/services"
	"github.com/joho/godotenv"

	"github.com/K44Z/kzchat/internal/server/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

func main() {
	envPath := filepath.Join(string(os.PathSeparator), "etc", "kzchat", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("error loading env file from :", envPath)
	}
	PORT := os.Getenv("PORT")
	config, err := configs.Load()
	if err != nil {
		log.Fatal("error loading env file from: ", envPath)
	}

	db, err := database.ConnectDb(config)
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	log.Println("Migrations applied")

	service := services.NewService(db)
	app := fiber.New()
	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
		},
	))

	app.Use(logger.New())
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello")
	})
	routes.SetupRoutes(app, service)
	err = app.Listen(PORT)
	if err != nil {
		log.Fatal(err)
	}
}
