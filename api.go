package main

import (
	"encoding/json"
	"io"
	"kzchat/server/schemas"
	"net/http"
	repository "kzchat/server/database/generated"

	tea "github.com/charmbracelet/bubbletea"
)

var API_URL = "http://localhost:4000"

func fetchMessages(config schemas.Config, recipientId string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest("GET", API_URL+"/message/recId/"+recipientId, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+config.Token)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var messages []repository.Message
		if err := json.Unmarshal(body, &messages); err != nil {
			return err
		}

		// return ChatFetchedMsg{Messages: messages}
		return nil
	}
}



