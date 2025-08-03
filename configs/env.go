package configs

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUrl string
	Port  string
}

func Load() (*Config, error) {
	envPath := filepath.Join(string(os.PathSeparator), "etc", "kzchat", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}
	return &Config{
		DbUrl: os.Getenv("DB_URL"),
		Port:  os.Getenv("PORT"),
	}, nil
}
