package scripts

import (
	"context"
	"fmt"
	"kzchat/server/database"
	repository "kzchat/server/database/generated"
	"kzchat/server/schemas"
	"kzchat/server/services"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var messages = []schemas.Message{
	{
		Content: "Hey, how are you?",
		Time:    time.Now().Add(-5 * time.Minute),
	},
	{
		Content: "I'm good, username! You?",
		Time:    time.Now().Add(-4 * time.Minute),
	},
	{
		Content: "Just working on a project.",
		Time:    time.Now().Add(-3 * time.Minute),
	},
	{
		Content: "Nice!",
		Time:    time.Now().Add(-2 * time.Minute),
	},
	{
		Content: "Yup! It's looking great so far.",
		Time:    time.Now().Add(-1 * time.Minute),
	},
}

var users = []repository.User{
	{
		Username: "username",
		ID:       1,
	},
	{
		Username: "amine",
		ID:       2,
	},
}

func SeedMessages() error {
	ctx := context.Background()
	chat, err := database.Queries.GetChatMessagesByChatId(ctx, 1)
	if err != nil {
		return err
	}
	fmt.Println(chat)
	if len(chat) == 0 {
		fmt.Println("seeding ...")
		chat, err := services.CreateDMChatFromMessage(messages[0], users)
		if err != nil {
			return err
		}
	
		for _ ,m := range messages {
			timestamp := pgtype.Timestamp{
				Time:  m.Time,
				Valid: true,
			}
			params := repository.StoreChatMessageParams{
				SenderID: users[0].ID,
				Content:  m.Content,
				ChatID:   chat.ID,
				Time:     timestamp,
				Type:     "dm",
			}
			_, err = database.Queries.StoreChatMessage(ctx, params)
			if err != nil {
				return fmt.Errorf("error storing message :%s", err)
			}
		}
	} else {
		fmt.Println("data already seeded")
		return nil
	}
	return nil
}
