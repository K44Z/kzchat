package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTtoken(user schemas.User) (string, error) {
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

func Authenticate(tokenString string) (*schemas.User, error) {
	JWT_SECRET := os.Getenv("JWT_SECRECT")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	user := schemas.User{
		ID:       int32(claims["sub"].(float64)),
		Username: claims["username"].(string),
	}

	return &user, nil
}
