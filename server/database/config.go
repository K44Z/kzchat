package database

import (
	"context"
	"log"
	"os"
	repository "kzchat/server/database/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Queries *repository.Queries

func ConnectDb() error {
	DB_URL := os.Getenv("DB_URL")
	dbConn, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatal("failed to connect to db :", err)
		return err
	}
	Queries = repository.New(dbConn)
	return nil
}
