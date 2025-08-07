package schemas

import (
	"time"
)

type Auth struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}
type Config struct {
	Username string `json:"username" validate:"required"`
	Token    string `json:"token" validate:"required"`
}

func (a Auth) Read(p []byte) (n int, err error) {
	panic("unimplemented")
}

type Message struct {
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
	Sender   User      `json:"sender"`
	Receiver User      `json:"receiver"`
	Chat     Chat      `json:"chat"`
}

type Chat struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type CreateMessageSchema struct {
	Content  string    `json:"content" validate:"required"`
	Type     string    `json:"type" validate:"required,oneof=dm chan"`
	SenderId int32     `json:"sender_id" validate:"required"`
	ChatId   int32     `json:"chat_id" validate:"required"`
	Time     time.Time `json:"time"`
}
type GetChatIdByParticipants struct {
	Members []string `json:"members" validate:"required"`
}

type User struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
}

type InternalUser struct {
	ID       int32
	Username string
	Password string
}

type CreateChatByMessage struct {
	Members []string `json:"members" validate:"required"`
	Message Message  `json:"message" validate:"required"`
}

func (u User) FilterValue() string {
	return u.Username
}

type ChatMember struct {
	ChatId int32 `json:"chatId" validate:"required"`
	UserId int32 `json:"userId" validate:"required"`
}
