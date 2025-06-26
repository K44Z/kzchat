package services

import (
	"context"
	"fmt"
	"kzchat/server/database"
	repository "kzchat/server/database/generated"
	"kzchat/server/schemas"
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
