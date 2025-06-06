package main

import (
	"kzchat/server/database"
	"kzchat/server/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	PORT := os.Getenv("PORT")
	err = database.ConnectDb()
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	log.Println("Migrations applied ✅")
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
	routes.SetupRoutes(app)
	err = app.Listen(PORT)
	if err != nil {
		log.Fatal(err)
	}
}
