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
	return &Services{
		UserService: NewUserService(repository.NewUserRepository(db)),
		ChatService: NewChatService(repository.NewChatRepository(db)),
	}
}

func wrap(err error, m string) error {
	if m == "" {
		return fmt.Errorf("Service : %w", err)
	}
	return fmt.Errorf("Handler : %s %w", m, err)
}
