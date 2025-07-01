package screens

import (
	"fmt"
	"kzchat/api"
	"kzchat/server/schemas"
	"os"

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
		if len(args) < 1 {
			return "Usage: dm <username>", nil
		}
		ctx.Model.Recipient.Username = args[0]

		return "", nil
	},
	"open": func(ctx CommandContext, args []string) (string, tea.Cmd) {
		if len(args) < 1 {
			return "Usage: open <username>", nil
		}

		id, users, err := api.GetChat([]string{api.Config.Username, args[0]})
		if users == nil {
			return err.Error(), nil
		}
		if err != nil {
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
