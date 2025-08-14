package services

import (
	"fmt"

	"github.com/K44Z/kzchat/internal/server/database"
	"github.com/K44Z/kzchat/internal/server/repository"
)

type Services struct {
	UserService UserService
	ChatService ChatService
}

func NewService(db *database.DB) *Services {
	userService := NewUserService(repository.NewUserRepository(db))
	return &Services{
		UserService: userService,
		ChatService: NewChatService(repository.NewChatRepository(db), userService),
	}
}

func wrap(err error, m string) error {
	if m == "" {
		return fmt.Errorf("Service : %w", err)
	}
	return fmt.Errorf("Handler : %s %w", m, err)
}
