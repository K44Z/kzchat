package database

import (
	"context"

	"github.com/K44Z/kzchat/configs"
	sqlc "github.com/K44Z/kzchat/internal/server/database/generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

var err error

type DB struct {
	DBTX sqlc.DBTX
	Pool *pgxpool.Pool
}

func ConnectDb(c *configs.Config) (*DB, error) {
	conn, err := pgxpool.New(context.Background(), c.DbUrl)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return &DB{
		DBTX: conn,
		Pool: conn,
	}, nil
}
