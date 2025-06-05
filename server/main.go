package main

import (
	"kzchat/server/database"
	"kzchat/server/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	PORT := os.Getenv("PORT")
	_, err = database.ConnectDb()
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
	routes.SetupRoutes(app)

	err =app.Listen(PORT)
	if err != nil {
		log.Fatal(err)
	}
}
