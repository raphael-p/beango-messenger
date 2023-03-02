package database

type Message struct {
	Id      string `json:"id"`
	ChatId  string `json:"chatId"`
	Content string `json:"content"`
}

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Key         []byte `json:"key"`
}
