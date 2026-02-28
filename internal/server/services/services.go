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
		return fmt.Errorf("service : %w", err)
	}
	return fmt.Errorf("handler : %s %w", m, err)
}

func CreateChatName(sender, receiver string) string {
	return fmt.Sprintf("%s - %s ", sender, receiver)
}
