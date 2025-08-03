package repository

import (
	"context"
	"fmt"

	"github.com/K44Z/kzchat/internal/server/database"
	sqlc "github.com/K44Z/kzchat/internal/server/database/generated"
	"github.com/K44Z/kzchat/internal/server/schemas"
)

type UserRepository interface {
	GetUsernameById(ctx context.Context, id int32) (*string, error)
	Create(ctx context.Context, user, password string) (*schemas.User, error)
	GetByUsername(ctx context.Context, username string) (*schemas.InternalUser, error)
	GetAll(ctx context.Context) ([]schemas.User, error)
}

type sqlcUserRepo struct {
	queries *sqlc.Queries
}

func NewUserRepository(db *database.DB) UserRepository {
	return &sqlcUserRepo{queries: sqlc.New(db.DBTX)}
}

func (u *sqlcUserRepo) GetUsernameById(ctx context.Context, id int32) (*string, error) {
	username, err := u.queries.GetUsernameById(ctx, id)
	if err != nil {
		return nil, wrap(err, "")
	}
	return &username, nil
}

func (u *sqlcUserRepo) Create(ctx context.Context, username, password string) (*schemas.User, error) {
	res, err := u.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, wrap(err, "")
	}
	return &schemas.User{
		ID:       res.ID,
		Username: res.Username,
	}, nil
}

func (u *sqlcUserRepo) GetByUsername(ctx context.Context, username string) (*schemas.InternalUser, error) {
	user, err := u.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, wrap(err, "")
	}
	return &schemas.InternalUser{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

func (u *sqlcUserRepo) GetAll(ctx context.Context) ([]schemas.User, error) {
	res, err := u.queries.GetUsers(ctx)
	if err != nil {
		return nil, wrap(err, "")
	}
	var users []schemas.User
	for _, user := range res {
		users = append(users, schemas.User{
			Username: user.Username,
			ID:       user.ID,
		})
	}
	return users, nil
}

func wrap(err error, m string) error {
	if m == "" {
		return fmt.Errorf("Repo : %w", err)
	}
	return fmt.Errorf("Repo : %s %w", m, err)
}
