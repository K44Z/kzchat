package screens

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/K44Z/kzchat/internal/api"

	tea "github.com/charmbracelet/bubbletea"
)

type CommandContext struct {
	Model *ChatModel
}

type CommandFunc func(ctx CommandContext, args []string) (string, tea.Cmd)

var CommandRegistry = map[string]CommandFunc{
	"q": func(ctx CommandContext, args []string) (string, tea.Cmd) {
		os.Exit(0)
		return "", tea.Quit
	},
	"clear": func(ctx CommandContext, args []string) (string, tea.Cmd) {
		ctx.Model.Messages = nil
		ctx.Model.Viewport.SetContent("")
		return "", nil
	},
	"dm": func(ctx CommandContext, args []string) (string, tea.Cmd) {
		if len(args) < 2 {
			return `Usage: dm <username> "<message>"`, nil
		}

		recipient := args[0]
		message := strings.Join(args[1:], " ")
		if message == "" {
			return "Message cannot be empty", nil
		}

		tempRecipient := ctx.Model.Recipient
		tempInput := ctx.Model.Textarea.Value()
		var (
			chatID *int32
			chat   schemas.Chat
		)

		chatID, users, err := api.GetChat([]string{api.Config.Username, recipient})
		if err != nil {
			var cusError *api.NotFoundErr
			if errors.As(err, &cusError) {
				m := schemas.Message{
					Content: message,
					Time:    time.Now(),
					Sender: schemas.User{
						Username: api.Config.Username,
					},
					Receiver: schemas.User{
						Username: recipient,
					},
				}
				chat, err = api.CreateChat(m)
				if err != nil {
					return err.Error(), nil
				}
				ctx.Model.Chat = chat
			} else {
				return err.Error(), nil
			}
		} else {
			ctx.Model.Chat.ID = *chatID
			ctx.Model.Recipient = users[1]
		}

		ctx.Model.Textarea.SetValue(message)
		ctx.Model.SendMessage()
		ctx.Model.Recipient = tempRecipient
		ctx.Model.Textarea.SetValue(tempInput)

		return "", nil
	},
	"open": func(ctx CommandContext, args []string) (string, tea.Cmd) {
		if len(args) < 1 || len(args) > 2 {
			return "Usage: open <username>", nil
		}
		id, users, err := api.GetChat([]string{api.Config.Username, args[0]})
		if id == nil || users == nil || err != nil {
			return fmt.Sprintf("Error opening chat: %v", err), nil
		}
		chat := schemas.Chat{
			Name: fmt.Sprint(api.Config.Username, " - ", args[0]),
			ID:   *id,
		}
		ctx.Model.Chat = chat
		ctx.Model.Current = users[0]
		ctx.Model.Recipient = users[1]

		// Return empty string with the command to ensure the command gets executed
		return "", ctx.Model.FetchMessages()
	},
}
