package services

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/K44Z/kzchat/internal/server/repository"
	"github.com/K44Z/kzchat/internal/server/schemas"
)

type UserService interface {
	CreateUser(ctx context.Context, username, password string) error
	GetUserByUsername(ctx context.Context, username string) (*schemas.User, error)
	CheckExistingUser(ctx context.Context, username string) (bool, error)
	GetUsernameById(ctx context.Context, id int32) (*string, error)
	GetAllUsersService(ctx context.Context) ([]schemas.User, error)
	GetUserWithPassword(ctx context.Context, username string) (*schemas.InternalUser, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(u repository.UserRepository) UserService {
	return &userService{
		userRepo: u,
	}
}

func (u *userService) CreateUser(ctx context.Context, username, password string) error {
	str, err := u.userRepo.Create(ctx, username, password)
	if err != nil {
		return wrap(err, "error creating user")
	}
	log.Println(str)
	return nil
}

func (u *userService) GetUserByUsername(ctx context.Context, username string) (*schemas.User, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &schemas.User{
		Username: user.Username,
		ID:       user.ID,
	}, nil

}

func (u *userService) CheckExistingUser(ctx context.Context, username string) (bool, error) {
	_, err := u.userRepo.GetByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, nil
	}
	return true, nil
}

func (u *userService) GetUsernameById(ctx context.Context, id int32) (*string, error) {
	username, err := u.userRepo.GetUsernameById(ctx, id)
	if err != nil {
		return nil, wrap(err, "error getting user by username")
	}
	return username, nil
}

func (u *userService) GetAllUsersService(ctx context.Context) ([]schemas.User, error) {
	var result []schemas.User
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, wrap(err, "error getting all users")
	}
	for _, user := range users {
		result = append(result, schemas.User{
			ID:       user.ID,
			Username: user.Username,
		})
	}
	return result, nil
}
func (u *userService) GetUserWithPassword(ctx context.Context, username string) (*schemas.InternalUser, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, wrap(err, "error getting user by username")
	}
	return &schemas.InternalUser{
		Username: user.Username,
		Password: user.Password,
		ID:       user.ID,
	}, nil
}
