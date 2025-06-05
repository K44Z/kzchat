package models

import "time"

type User struct {
	Id        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
}

type Message struct {
	SenderId   string
	ReceiverId string 
	Content    string
	Time       time.Time
	Type       string 
}

type Chat struct {
	Messages []Message
	Members  []User
}


type Config struct {
	Username string 
	Token string
}