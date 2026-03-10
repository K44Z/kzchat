package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/K44Z/kzchat/internal/server/database"
	sqlc "github.com/K44Z/kzchat/internal/server/database/generated"
	"github.com/K44Z/kzchat/internal/server/schemas"
)

type ChatRepository interface {
	Create(ctx context.Context, arg sqlc.CreateChatParams, users []schemas.User) (*schemas.Chat, error)
	CreateMembers(ctx context.Context, arg sqlc.CreateChatMembersParams) error
	FindByParticipants(ctx context.Context, arg []int32) (*int32, string, error)
	GetById(ctx context.Context, id int32) (*schemas.Chat, error)
	GetMessagesByChatId(ctx context.Context, id int32) ([]schemas.Message, error)
	GetMessagesByParticipants(ctx context.Context, arg sqlc.GetDmChatMessagesByParticipantsParams) ([]schemas.Message, error)
	StoreMessage(ctx context.Context, arg sqlc.StoreChatMessageParams) error
	// GetAttachementByMessage(ctx context.Context, id int32) ([]schemas.Attachment, error)
	GetAttachementsByChat(ctx context.Context, id int32) ([]schemas.Attachment, error)
	SaveAttachmentRepo(ctx context.Context, param schemas.SaveAttachementParams) error
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

func (c *chatRepository) Create(ctx context.Context, arg sqlc.CreateChatParams, users []schemas.User) (*schemas.Chat, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %s", err)
	}
	defer tx.Rollback(ctx)

	ttx := c.queries.WithTx(tx)

	chat, err := ttx.CreateChat(ctx, arg)
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
		Name: arg.Name,
	}, nil
}

func (c *chatRepository) CreateMembers(ctx context.Context, arg sqlc.CreateChatMembersParams) error {
	res, err := c.queries.CreateChatMembers(ctx, arg)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("no Rows affected")
	}
	return nil
}

func (c *chatRepository) FindByParticipants(ctx context.Context, arg []int32) (*int32, string, error) {
	res, err := c.queries.FindChatByParticipants(ctx, arg)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", fiber.ErrNotFound
	}
	if err != nil {
		return nil, "", err
	}
	return &res.ChatID, res.ChatName, nil
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
	fmt.Println("repo ", messages)
	return messages, nil
}

func (c *chatRepository) StoreMessage(ctx context.Context, arg sqlc.StoreChatMessageParams) error {
	res, err := c.queries.StoreChatMessage(ctx, arg)
	if err != nil {
		return err
	}
	if count := res.RowsAffected(); count == 0 {
		return fmt.Errorf("no Rows affected")
	}
	return nil
}

// func (c *chatRepository) GetAttachementByMessage(ctx context.Context, id int32) ([]schemas.Attachment, error) {
// 	var result []schemas.Attachment
// 	res, err := c.queries.GetAttachementsByMessageId(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, a := range res {
// 		result = append(result, schemas.Attachment{
// 			ID:       a.ID,
// 			FileName: a.FileName,
// 			FileType: a.FileType,
// 			FileSize: a.FileSize,
// 			URL:      a.FileUrl,
// 		})
// 	}
// 	return result, nil
// }

func (c *chatRepository) GetAttachementsByChat(ctx context.Context, id int32) ([]schemas.Attachment, error) {
	var result []schemas.Attachment
	res, err := c.queries.GetAttachementsByChatId(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, a := range res {
		result = append(result, schemas.Attachment{
			ID:       a.ID,
			FileName: a.FileName,
			FileType: a.FileType,
			FileSize: a.FileSize,
			URL:      a.FileUrl,
		})
	}
	return result, nil
}

func (c *chatRepository) SaveAttachmentRepo(ctx context.Context, param schemas.SaveAttachementParams) error {
	_, err := c.queries.SaveAttachment(ctx, sqlc.SaveAttachmentParams{
		FileName: param.File.Filename,
		FileType: param.File.Header["Content-Type"][0],
		FileSize: int32(param.File.Size),
		ChatID:   param.ChatID,
		FileUrl:  param.Path,
	})
	if err != nil {
		return err
	}
	return nil
}
