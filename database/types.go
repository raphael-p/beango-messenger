package database

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
	Id       string `json:"id"`
	Username string `json:"username"`
	Key      string `json:"key"`
}

type ChatObject interface {
	Message | Chat | User
}
