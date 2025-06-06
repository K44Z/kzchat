package services

import (
	"context"
	"database/sql"
	"errors"
	"kzchat/server/database"
	repository "kzchat/server/database/generated"
	"log"
)

func CreateUser(user repository.CreateUserParams) error {
	ctx := context.Background()
	str, err := database.Queries.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	log.Println(str)
	return nil
}

func GetUserByUsername(username string) (repository.User, error) {
	ctx := context.Background()
	user, err := database.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return repository.User{}, err
	}
	return user, nil

}

func CheckExistingUser(username string) (bool, error) {
	ctx := context.Background()
	_, err := database.Queries.GetUserByUsername(ctx, username)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, nil
	}
	return true, nil
}
