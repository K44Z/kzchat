package screens

import (
	"errors"
	"fmt"
	"kzchat/api"
	"kzchat/server/schemas"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type CommandContext struct {
	Model *ChatModel
}

type CommandFunc func(ctx CommandContext, args []string) (string, tea.Cmd)

var CommandRegistry = map[string]CommandFunc{
	"quit": func(ctx CommandContext, args []string) (string, tea.Cmd) {
		os.Exit(0)
		return "", nil
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
		tempInput := ctx.Model.Input.Value()
		var (
			chatID int32
			chat schemas.Chat
		)

		chatID, users, err := api.GetChat([]string{api.Config.Username, recipient})
		if err != nil {
			var cusError *api.NotFoundErr
			if errors.As(err, &cusError){
				m := schemas.Message{
					Content: message,
					Time: time.Now(),
					SenderUsername: api.Config.Username,
					ReceiverUsername: recipient,
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
			ctx.Model.Chat.ID = chatID
			ctx.Model.Recipient = users[1]
		}

		ctx.Model.Input.SetValue(message)
		ctx.Model.SendMessage()
		ctx.Model.Recipient = tempRecipient
		ctx.Model.Input.SetValue(tempInput)

		return "", nil
	},
	"open": func(ctx CommandContext, args []string) (string, tea.Cmd) {
		if len(args) < 1 {
			return "Usage: open <username>", nil
		}

		id, users, err := api.GetChat([]string{api.Config.Username, args[0]})
		if users == nil || err != nil {
			return err.Error(), nil
		}
		chat := schemas.Chat{
			Name: fmt.Sprint(api.Config.Username, " - ", args[0]),
			ID:   id,
		}
		ctx.Model.Chat = chat
		ctx.Model.Recipient.Username = args[0]
		cmd := ctx.Model.FetchMessages()
		ctx.Model.Current = users[0]
		ctx.Model.Recipient = users[1]
		return "", cmd
	},
}
