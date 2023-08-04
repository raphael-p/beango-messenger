package database

import (
	"time"
)

type Message struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"userID"`
	ChatID        int64     `json:"chatID"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

type MessageExtended struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"userID"`
	ChatID          int64     `json:"chatID"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"createdAt"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt"`
	UserDisplayName string    `json:"userDisplayName"`
}

func (conn *MongoConnection) GetMessagesByChatID(chatID int64) ([]MessageExtended, error) {
	return scanRows[MessageExtended](conn.Query(
		`SELECT
			m.*,
			u.display_name as user_display_name
		FROM message m
		LEFT JOIN "user" u ON u.id = m.user_id
		WHERE chat_id = $1
		ORDER BY m.created_at ASC;`,
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
