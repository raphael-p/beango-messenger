package database

import "time"

type Message struct {
	ID            int       `json:"id"`
	UserID        int       `json:"userID"`
	ChatID        int       `json:"chatID"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"LastUpdatedAt"`
}

func (conn *MongoConnection) GetMessagesByChatID(chatID int) []Message {
	messages := []Message{}
	for _, message := range Messages {
		if message.ChatID == chatID {
			messages = append(messages, message)
		}
	}
	return messages
}

func (conn *MongoConnection) SetMessage(message *Message) {
	Messages[message.ID] = *message
}
