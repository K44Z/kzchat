package database

import (
	"kzchat/server/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func ConnectDb() (*gorm.DB, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}
	DB_URL := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(DB_URL), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database ", err)
		return nil, err
	}
	err = Migrate()
	if err != nil {
		log.Fatal("Error applying migrations")
		return nil, err
	}
	log.Println("Database Connected ✅")
	return DB, nil
}

func Migrate() error {
	err = DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return err
	}
	return nil
}
