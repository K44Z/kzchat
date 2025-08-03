package services

import (
	"context"
	"errors"
	"fmt"

	sqlc "github.com/K44Z/kzchat/internal/server/database/generated"
	"github.com/K44Z/kzchat/internal/server/repository"
	"github.com/K44Z/kzchat/internal/server/schemas"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ChatService interface {
	CreateChatFromMessage(ctx context.Context, m schemas.Message, users []schemas.User) (*schemas.Chat, error)
	GetChatById(ctx context.Context, id int32) (*schemas.Chat, error)
	GetMessagesByParticipants(ctx context.Context, current, rec schemas.User) ([]schemas.Message, error)
	CreateDM(ctx context.Context, m schemas.Message) error
	GetChatIdByParticipants(ctx context.Context, arg sqlc.FindChatByParticipantsParams) (*int32, error)
}

type chatService struct {
	chatRepo    repository.ChatRepository
	userService UserService
}

func NewChatService(c repository.ChatRepository) ChatService {
	return &chatService{
		chatRepo: c,
	}
}

func (c *chatService) CreateChatFromMessage(ctx context.Context, m schemas.Message, users []schemas.User) (*schemas.Chat, error) {
	name := fmt.Sprintf("%s - %s ", m.Sender.Username, m.Receiver.Username)
	fmt.Println(name)
	chat, err := c.chatRepo.Create(ctx, sqlc.CreateChatParams{
		Type: "dm",
		Name: name,
	}, users, name)
	if err != nil {
		return nil, wrap(err, "")
	}
	return chat, nil
}

func (c *chatService) GetChatById(ctx context.Context, id int32) (*schemas.Chat, error) {
	chat, err := c.chatRepo.GetById(ctx, id)
	if err != nil {
		return nil, wrap(err, "")
	}
	return chat, nil
}

func (c *chatService) GetMessagesByParticipants(ctx context.Context, current, rec schemas.User) ([]schemas.Message, error) {
	messages, err := c.chatRepo.GetMessagesByParticipants(ctx, sqlc.GetDmChatMessagesByParticipantsParams{
		UserID:   current.ID,
		UserID_2: rec.ID,
	})
	if err != nil {
		return nil, wrap(err, "")
	}
	return messages, nil
}

func (c *chatService) CreateDM(ctx context.Context, m schemas.Message) error {
	var (
		chat      *schemas.Chat
		err       error
		usernames = []string{
			m.Sender.Username, m.Receiver.Username,
		}
		users []schemas.User
	)

	for _, u := range usernames {
		user, err := c.userService.GetUserByUsername(ctx, u)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return wrap(err, "user does not exist")
			} else {
				return wrap(err, "")
			}
		}
		users = append(users, *user)
	}
	fmt.Println("the chat id is :", m.Chat.ID)
	if m.Chat.ID == 0 {
		chat, err = c.CreateChatFromMessage(ctx, m, users)
		if err != nil {
			return wrap(err, "error creating chat")
		}
	} else {
		chat = &schemas.Chat{ID: m.Chat.ID}
	}

	timestamp := pgtype.Timestamp{
		Time:  m.Time,
		Valid: true,
	}
	params := sqlc.StoreChatMessageParams{
		SenderID: users[0].ID,
		Content:  m.Content,
		ChatID:   chat.ID,
		Time:     timestamp,
		Type:     "dm",
	}
	err = c.chatRepo.StoreMessage(ctx, params)
	if err != nil {
		return wrap(err, "error storing message")
	}
	return nil
}

func (s *chatService) GetChatIdByParticipants(ctx context.Context, arg sqlc.FindChatByParticipantsParams) (id *int32, err error) {
	defer wrap(err, "")
	return s.chatRepo.FindByParticipants(ctx, arg)
}
