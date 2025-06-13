package schemas

import "time"

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

type CreateMessageSchema struct {
	Content  string    `json:"content" validate:"required"`
	Type     string    `json:"type" validate:"required,oneof=dm chan"`
	SenderId int32     `json:"sender_id" validate:"required"`
	ChatId   int32     `json:"chat_id" validate:"required"`
	Time     time.Time `json:"time"`
}
