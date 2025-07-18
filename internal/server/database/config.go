package database

import (
	"context"
	"log"
	"os"

	repository "github.com/K44Z/kzchat/internal/server/database/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Queries *repository.Queries
var DbConn *pgxpool.Pool
var err error

func ConnectDb() error {
	DB_URL := os.Getenv("DB_URL")
	DbConn, err = pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatal("failed to connect to db :", err)
		return err
	}
	err = DbConn.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
	Queries = repository.New(DbConn)
	return nil
}
