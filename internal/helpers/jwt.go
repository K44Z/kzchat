package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTtoken(user schemas.User) (string, error) {
	JWT_SECRET := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	})
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Authenticate(tokenString string) (*schemas.User, error) {
	JWT_SECRET := os.Getenv("JWT_SECRET")
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

	sub, ok := claims["sub"]
	if !ok {
		return nil, fmt.Errorf("missing sub claim")
	}
	username, ok := claims["username"]
	if !ok {
		return nil, fmt.Errorf("missing username claim")
	}
	subFloat, ok := sub.(float64)
	if !ok {
		return nil, fmt.Errorf("invalid sub claim type")
	}
	usernameStr, ok := username.(string)
	if !ok {
		return nil, fmt.Errorf("invalid username claim type")
	}

	user := schemas.User{
		ID:       int32(subFloat),
		Username: usernameStr,
	}

	return &user, nil
}
