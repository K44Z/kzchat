package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/K44Z/kzchat/internal/server/database"
	sqlc "github.com/K44Z/kzchat/internal/server/database/generated"
	"github.com/K44Z/kzchat/internal/server/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository interface {
	Create(ctx context.Context, arg sqlc.CreateChatParams, users []schemas.User, name string) (*schemas.Chat, error)
	CreateMembers(ctx context.Context, arg sqlc.CreateChatMembersParams) error
	FindByParticipants(ctx context.Context, arg []int32) (*int32, error)
	GetById(ctx context.Context, id int32) (*schemas.Chat, error)
	GetMessagesByChatId(ctx context.Context, id int32) ([]schemas.Message, error)
	GetMessagesByParticipants(ctx context.Context, arg sqlc.GetDmChatMessagesByParticipantsParams) ([]schemas.Message, error)
	StoreMessage(ctx context.Context, arg sqlc.StoreChatMessageParams) error
}

type chatRepository struct {
	queries sqlc.Queries
	db      *pgxpool.Pool
}

func NewChatRepository(db *database.DB) ChatRepository {
	return &chatRepository{
		queries: *sqlc.New(db.DBTX),
		db:      db.Pool,
	}
}

func (c *chatRepository) Create(ctx context.Context, arg sqlc.CreateChatParams, users []schemas.User, name string) (*schemas.Chat, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %s", err)
	}
	defer tx.Rollback(ctx)

	ttx := c.queries.WithTx(tx)

	chat, err := ttx.CreateChat(ctx, sqlc.CreateChatParams{
		Type: "dm",
		Name: name,
	})
	if err != nil {
		return nil, wrap(err, "")
	}

	for _, u := range users {
		_, err := ttx.CreateChatMembers(ctx, sqlc.CreateChatMembersParams{
			ChatID: chat.ID,
			UserID: u.ID,
		})
		if err != nil {
			return nil, wrap(err, "")
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, wrap(err, "")
	}

	return &schemas.Chat{
		ID:   chat.ID,
		Name: name,
	}, nil
}

func (c *chatRepository) CreateMembers(ctx context.Context, arg sqlc.CreateChatMembersParams) error {
	res, err := c.queries.CreateChatMembers(ctx, arg)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("No Rows affected")
	}
	return nil
}

func (c *chatRepository) FindByParticipants(ctx context.Context, arg []int32) (*int32, error) {
	id, err := c.queries.FindChatByParticipants(ctx, arg)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fiber.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (c *chatRepository) GetById(ctx context.Context, id int32) (*schemas.Chat, error) {
	chat, err := c.queries.GetChatById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fiber.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &schemas.Chat{
		ID:   chat.ID,
		Name: chat.Name,
	}, nil
}

func (c *chatRepository) GetMessagesByChatId(ctx context.Context, id int32) ([]schemas.Message, error) {
	res, err := c.queries.GetChatMessagesByChatId(ctx, id)
	if err != nil {
		return nil, err
	}
	var messages []schemas.Message
	for _, message := range res {
		messages = append(messages, schemas.Message{
			Content: message.Content,
			Time:    message.Time.Time,
		})
	}
	return messages, nil
}

func (c *chatRepository) GetMessagesByParticipants(ctx context.Context, arg sqlc.GetDmChatMessagesByParticipantsParams) ([]schemas.Message, error) {
	res, err := c.queries.GetDmChatMessagesByParticipants(ctx, arg)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fiber.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var messages []schemas.Message
	for _, message := range res {
		sender := schemas.User{
			ID:       message.SenderID,
			Username: message.SenderUsername,
		}
		receiver := schemas.User{
			ID:       message.ReceiverID,
			Username: message.ReceiverUsername,
		}
		messages = append(messages, schemas.Message{
			Content:  message.Content,
			Time:     message.Time.Time,
			Sender:   sender,
			Receiver: receiver,
		})
	}
	return messages, nil
}

func (c *chatRepository) StoreMessage(ctx context.Context, arg sqlc.StoreChatMessageParams) error {
	res, err := c.queries.StoreChatMessage(ctx, arg)
	if err != nil {
		return err
	}
	if count := res.RowsAffected(); count == 0 {
		return fmt.Errorf("No Rows affected")
	}
	return nil
}
