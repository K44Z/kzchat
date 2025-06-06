package authentication

import (
	"os"
	"time"

	repository "kzchat/server/database/generated"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateJWTtoken(user repository.User) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return err.Error(), err
	}

	JWT_SECRET := os.Getenv("JWT_SECRECT")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	})
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return err.Error(), err
	}
	return tokenString, nil
}

func Authenticate(tokenString string) (repository.User, error) {
	return repository.User{}, nil
}
