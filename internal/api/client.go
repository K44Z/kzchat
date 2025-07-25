package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/K44Z/kzchat/internal/server/schemas"

	"github.com/golang-jwt/jwt/v5"
)

var BASE_URL string
var WS_URL string

type Claims struct {
	Username string `json:"username"`
	Sub      string `json:"sub"`
	jwt.RegisteredClaims
}

type GetChatResponse struct {
	ChatId int32          `json:"chatId"`
	Users  []schemas.User `json:"users"`
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
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, "/.config/kzchat")

	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		if err = os.Mkdir(configDir, 0700); err != nil {
			return err
		}
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	tokenFile := filepath.Join(configDir, "token.json")
	return ioutil.WriteFile(tokenFile, data, 0600)
}

func ReadConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(filepath.Join(home, "/.config/kzchat", "token.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &Config); err != nil {
		return err
	}
	return nil
}

func IsTokenValid(tokenString string) bool {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected singing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return false
	}

	return claims.ExpiresAt.After(time.Now())
}

func GetChat(m []string) (int32, []schemas.User, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(map[string]any{
		"members": m,
	})
	if err != nil {
		return 0, nil, fmt.Errorf("error marshaling data: %w", err)
	}
	url := fmt.Sprintf("%s/messages/chat", BASE_URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Config.Token)
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return 0, nil, &NotFoundErr{
				Msg: `no previous chat was found, use dm <username> <"message">`,
			}
		}
		return 0, nil, fmt.Errorf("unexpected status code %d ", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("error reading response: %w", err)
	}

	var result GetChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return result.ChatId, result.Users, nil
}

func CreateChat(message schemas.Message) (schemas.Chat, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(map[string]any{
		"members": []string{message.SenderUsername, message.ReceiverUsername},
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
