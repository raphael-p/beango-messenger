package database

import "time"

type Message struct {
	Id      string `json:"id"`
	ChatId  string `json:"chatId"`
	Content string `json:"content"`
}

type Chat struct {
	Id      string   `json:"id"`
	UserIds []string `json:"userIds"`
}

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Key         []byte `json:"key"`
}

type Session struct {
	Id         string    `json:"id"`
	UserId     string    `json:"userId"`
	ExpiryDate time.Time `json:"expiryDate"`
}
