package authentication

import (
	"kzchat/server/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateJWTtoken(user models.User) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return err.Error(), err
	}

	JWT_SECRET := os.Getenv("JWT_SECRECT")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.Id,
		"username": user.Username,
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	})
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return err.Error(), err
	}
	return tokenString, nil
}

func Authenticate(tokenString string) (models.User, error) {
	return models.User{}, nil
}
