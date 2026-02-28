package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/list"

	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/golang-jwt/jwt/v5"
)

var (
	BASE_URL string
	WS_URL   string
)

type ApiResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

type GetChatResponse struct {
	Status string   `json:"status"`
	Data   ChatData `json:"data"`
}

type ChatData struct {
	ChatId int32          `json:"chatId"`
	Users  []schemas.User `json:"users"`
}

type FetchedUsersReponse struct {
	Status string          `json:"status"`
	Data   UserListReponse `json:"data"`
}

type UserListReponse struct {
	Users []schemas.User `json:"users"`
}

type Claims struct {
	Username string `json:"username"`
	Sub      string `json:"sub"`
	jwt.RegisteredClaims
}

type CreateChatResponse struct {
	Chat schemas.Chat `json:"chat"`
}

type NotFoundErr struct {
	Msg string
}

func (e *NotFoundErr) Error() string {
	return fmt.Sprint(e.Msg)
}

var Config schemas.Config

func SaveConfig(config schemas.Config) error {
	Config.Token = config.Token
	Config.Username = config.Username
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, "/.config/kzchat")

	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		if err = os.Mkdir(configDir, 0o700); err != nil {
			return err
		}
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	tokenFile := filepath.Join(configDir, "token.json")
	return ioutil.WriteFile(tokenFile, data, 0o600)
}

func ReadConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filepath.Join(home, "/.config/kzchat", "token.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &Config); err != nil {
		return err
	}
	return nil
}

func IsTokenValid(tokenString string) bool {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return false
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return false
	}
	return claims.ExpiresAt.After(time.Now())
}

func GetChat(m []string) (*int32, []schemas.User, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(map[string]any{
		"members": m,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling data: %w", err)
	}
	url := fmt.Sprintf("%s/messages/chat", BASE_URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Config.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil, &NotFoundErr{
				Msg: `no previous chat was found, use dm <username> <"message">`,
			}
		}
		return nil, nil, fmt.Errorf("unexpected status code %d ", resp.StatusCode)
	}

	var apiResp GetChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, nil, err
	}
	return &apiResp.Data.ChatId, apiResp.Data.Users, nil
}

func CreateChat(message schemas.Message) (schemas.Chat, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(map[string]any{
		"members": []string{message.Sender.Username, message.Receiver.Username},
		"message": message,
	})
	if err != nil {
		return schemas.Chat{}, fmt.Errorf("error marshaling data: %w", err)
	}
	url := fmt.Sprintf("%s/messages/createChat", BASE_URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return schemas.Chat{}, fmt.Errorf("error creating request : %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Config.Token)
	resp, err := client.Do(req)
	if err != nil {
		return schemas.Chat{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusBadRequest {
			return schemas.Chat{}, fmt.Errorf(`no previous chat was found, use dm <username> <"message">`)
		}
		return schemas.Chat{}, fmt.Errorf("unexpected status code %d ", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return schemas.Chat{}, fmt.Errorf("error unmarshaling response: %w", err)
	}
	var result CreateChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return schemas.Chat{}, fmt.Errorf("error unmarshaling response: %w", err)
	}
	return result.Chat, nil
}

func GetUsers() ([]list.Item, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", BASE_URL+"/users/usernames/all", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Config.Token)
	response, err := client.Do(req)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("unexpected status code :%d", response.StatusCode)
		}
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, err
	}
	defer response.Body.Close()
	var res FetchedUsersReponse
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	var list []list.Item
	for _, v := range res.Data.Users {
		list = append(list, schemas.User{
			Username: v.Username,
			ID:       v.ID,
		})
	}
	return list, nil
}
