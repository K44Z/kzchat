package services

import (
	"context"
	"fmt"

	repository "github.com/K44Z/kzchat/internal/server/database/generated"
	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/K44Z/kzchat/internal/server/database"
)

func CreateDMChatFromMessage(m schemas.Message, users []repository.User) (*repository.Chat, error) {
	ctx := context.Background()
	name := fmt.Sprintf("%s - %s ", m.SenderUsername, m.ReceiverUsername)
	fmt.Println(name)

	tx, err := database.DbConn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %s", err)
	}
	defer tx.Rollback(ctx)

	ttx := database.Queries.WithTx(tx)

	chatParams := repository.CreateChatParams{
		Type: "dm",
		Name: name,
	}

	chat, err := ttx.CreateChat(ctx, chatParams)
	if err != nil {
		return nil, fmt.Errorf("error creating chat, internal server error: %s", err)
	}

	for _, u := range users {
		_, err := ttx.CreateChatMembers(ctx, repository.CreateChatMembersParams{
			ChatID: chat.ID,
			UserID: u.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("error creating chat members %s", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %s", err)
	}
	return &chat, nil
}

func GetChatByIdService(id int32) (repository.Chat, error) {
	ctx := context.Background()
	chat, err := database.Queries.GetChatById(ctx, id)
	if err != nil {
		return repository.Chat{}, err
	}
	return chat, nil
}

func MapMessagesToClient(arr []repository.Message, users []repository.User) ([]schemas.Message, error) {
	var messages []schemas.Message
	for _, v := range arr {
		sender, err := GetUsernameById(v.SenderID)
		if err != nil {
			return nil, err
		}
		c, err := GetChatByIdService(v.ChatID)
		if err != nil {
			return nil, err
		}
		chat := MapChatToClient(c)
		ms := schemas.Message{
			Content:        v.Content,
			Time:           v.Time.Time,
			SenderUsername: sender,
			Chat:           chat,
		}
		messages = append(messages, ms)
	}
	return messages, nil
}

func MapChatToClient(chat repository.Chat) schemas.Chat {
	return schemas.Chat{
		ID:   chat.ID,
		Name: chat.Name,
	}
}
