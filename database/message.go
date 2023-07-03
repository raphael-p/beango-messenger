package database

import "time"

type Message struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"userID"`
	ChatID        int64     `json:"chatID"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"LastUpdatedAt"`
}

func (conn *MongoConnection) GetMessagesByChatID(chatID int64) []Message {
	messages := []Message{}
	for _, message := range Messages {
		if message.ChatID == chatID {
			messages = append(messages, message)
		}
	}
	return messages
}

func (conn *MongoConnection) SetMessage(message *Message) *Message {
	message.ID = int64(len(Messages) + 1)
	Messages[message.ID] = *message
	return message
}
