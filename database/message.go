package database

import (
	"time"
)

type MessageDatabase struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"userID"`
	ChatID        int64     `json:"chatID"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

type Message struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"userID"`
	ChatID          int64     `json:"chatID"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"createdAt"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt"`
	UserDisplayName string    `json:"userDisplayName"`
}

func (conn *MongoConnection) GetMessagesByChatID(chatID int64, fromMessageID int64) ([]Message, error) {
	return scanRows[Message](conn.Query(
		`SELECT
			m.*,
			u.display_name as user_display_name
		FROM message m
		LEFT JOIN "user" u ON u.id = m.user_id
		WHERE chat_id = $1 AND m.id > $2
		ORDER BY m.created_at ASC;`,
		chatID, fromMessageID,
	))
}

func (conn *MongoConnection) SetMessage(message *MessageDatabase) (*MessageDatabase, error) {
	return scanRow[MessageDatabase](conn.QueryRow(
		`INSERT INTO message (user_id, chat_id, content)
		VALUES ($1, $2, $3)
		RETURNING *`,
		message.UserID, message.ChatID, message.Content,
	))
}
