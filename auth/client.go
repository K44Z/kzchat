package authentication

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"kzchat/server/schemas"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	Username string `json:"username"`
	Sub      string `json:"sub"`
	jwt.RegisteredClaims
}

var Config schemas.Config

func SaveConfig(config schemas.Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, ".kzchat")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0700); err != nil {
			return err
		}
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	tokenFile := filepath.Join(configDir, "token.json")
	return ioutil.WriteFile(tokenFile, data, 0600)
}

func ReadConfig() (error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return  err
	}

	data, err := ioutil.ReadFile(filepath.Join(home, ".kzchat", "token.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &Config); err != nil {
		return  err
	}
	return nil
}

func IsTokenValid(tokenString string) bool {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected singing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return false
	}

	return claims.ExpiresAt.After(time.Now())
}
