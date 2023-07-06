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

func (conn *MongoConnection) GetMessagesByChatID(chatID int64) ([]Message, error) {
	return scanRows[Message](conn.Query(
		`SELECT * FROM message WHERE chat_id = $1
		ORDER BY created_at DESC`,
		chatID,
	))
}

func (conn *MongoConnection) SetMessage(message *Message) (*Message, error) {
	return scanRow[Message](conn.QueryRow(
		`INSERT INTO message (user_id, chat_id, content)
		VALUES ($1, $2, $3)
		RETURNING *`,
		message.UserID, message.ChatID, message.Content,
	))
}
