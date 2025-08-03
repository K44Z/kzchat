package scripts

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/K44Z/kzchat/internal/server/schemas"
// 	"github.com/K44Z/kzchat/internal/server/services"
// )

// var messages = []schemas.Message{
// 	{
// 		Content: "Hey, how are you?",
// 		Time:    time.Now().Add(-5 * time.Minute),
// 	},
// 	{
// 		Content: "I'm good, username! You?",
// 		Time:    time.Now().Add(-4 * time.Minute),
// 	},
// 	{
// 		Content: "Just working on a project.",
// 		Time:    time.Now().Add(-3 * time.Minute),
// 	},
// 	{
// 		Content: "Nice!",
// 		Time:    time.Now().Add(-2 * time.Minute),
// 	},
// 	{
// 		Content: "Yup! It's looking great so far.",
// 		Time:    time.Now().Add(-1 * time.Minute),
// 	},
// }

// var users = []schemas.User{
// 	{
// 		Username: "username",
// 		ID:       1,
// 	},
// 	{
// 		Username: "amine",
// 		ID:       2,
// 	},
// 	{
// 		Username: "username",
// 		ID:       8,
// 	},
// 	{
// 		Username: "amine",
// 		ID:       9,
// 	},
// 	{
// 		Username: "alex",
// 		ID:       10,
// 	},
// 	{
// 		Username: "maria",
// 		ID:       11,
// 	},
// 	{
// 		Username: "sdfksd",
// 		ID:       12,
// 	},
// 	{
// 		Username: "ksjfk",
// 		ID:       13,
// 	}, {
// 		Username: "john",
// 		ID:       14,
// 	}, {
// 		Username: "hij",
// 		ID:       15,
// 	},
// }

// func SeedMessages(s *services.Services) error {
// 	ctx := context.Background()
// 	chat, err := s.ChatService.GetChatById(ctx, 1)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(chat)
// 	if chat.ID == 0 {
// 		fmt.Println("seeding ...")
// 		chat, err := services.ChatService.CreateChatFromMessage(ctx, messages, users)
// 		if err != nil {
// 			return err
// 		}

// 		for _, m := range messages {
// 			// timestamp := pgtype.Timestamp{
// 			// 	Time:  m.Time,
// 			// 	Valid: true,
// 			// }
// 			// params := sqlc.StoreChatMessageParams{
// 			// 	SenderID: users[0].ID,
// 			// 	Content:  m.Content,
// 			// 	ChatID:   chat.ID,
// 			// 	Time:     timestamp,
// 			// 	Type:     "dm",
// 			// }
// 			err = s.ChatService.CreateDM(ctx, m)
// 			if err != nil {
// 				return fmt.Errorf("error storing message :%s", err)
// 			}
// 		}
// 	} else {
// 		fmt.Println("data already seeded")
// 		return nil
// 	}
// 	return nil
// }
