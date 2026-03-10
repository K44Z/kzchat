package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

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
	ChatId      int32                `json:"chatId"`
	Users       []schemas.User       `json:"users"`
	Name        string               `json:"name"`
	Attachments []schemas.Attachment `json:"attachments,omitempty"`
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

func IsTokenValid(tokenString string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", BASE_URL+"/auth/validate", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenString)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("token is invalid")
	}
	return nil
}

func GetChat(m []string) (*schemas.Chat, []schemas.User, error) {
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
	return &schemas.Chat{
		ID:          apiResp.Data.ChatId,
		Name:        apiResp.Data.Name,
		Attachments: apiResp.Data.Attachments,
	}, apiResp.Data.Users, nil
}

func CreateChat(message schemas.Message) (*schemas.Chat, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(map[string]any{
		"members": []string{message.Sender.Username, message.Receiver.Username},
		"message": message,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %w", err)
	}
	url := fmt.Sprintf("%s/messages/createChat", BASE_URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request : %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Config.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf(`no previous chat was found, use dm <username> <"message">`)
		}
		return nil, fmt.Errorf("unexpected status code %d ", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}
	var result CreateChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}
	return &result.Chat, nil
}

func GetUsers() ([]list.Item, error) {
	client := &http.Client{}
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

func UploadFile(path string, chatID int32) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, err := w.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return fmt.Errorf("failed to create form file %w", err)
	}
	if _, err = io.Copy(part, f); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	w.Close()
	id := strconv.Itoa(int(chatID))
	url := fmt.Sprintf("/messages/upload/chat/%v", id)
	req, err := http.NewRequest("POST", BASE_URL+url, &b)
	if err != nil {
		return fmt.Errorf("failed to create request %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Add("Authorization", "Bearer "+Config.Token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

func DownloadFile(chatID int32, name string) (tea.Cmd, error) {
	id := strconv.Itoa(int(chatID))
	url := BASE_URL + "/messages/chat/" + id + "/downloadFile"
	client := &http.Client{}
	b := struct {
		Path string `json:"path"`
	}{
		Path: "./uploads/" + name,
	}
	body, err := json.Marshal(&b)
	if err != nil {
		return nil, fmt.Errorf("Error marshing body %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Error creating request %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Config.Token)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %v", res.StatusCode)
	}
	defer res.Body.Close()
	err = os.MkdirAll("./downloads", os.ModePerm)
	if err != nil {
		return nil, err
	}
	out, err := os.Create("./downloads/" + name)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}
	return func() tea.Msg {
		return fmt.Sprintf("file %s downloaded successfully", name)
	}, nil
}
