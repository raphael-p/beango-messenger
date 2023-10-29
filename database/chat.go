package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type chatType string

// TODO: add support for NOTE and GROUP CHAT
const (
	NOTE         chatType = "note"
	PRIVATE_CHAT chatType = "private"
	GROUP_CHAT   chatType = "group"
)

type Chat struct {
	ID            int64     `json:"id"`
	Type          chatType  `json:"type"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

type ChatUser struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chatID"`
	UserID    int64     `json:"userID"`
	CreatedAt time.Time `json:"createdAt"`
}

func (conn *MongoConnection) GetChat(id, userID int64) (*Chat, error) {
	row := conn.QueryRow(
		`SELECT * FROM chat
		WHERE id = $1 
		AND EXISTS (
			SELECT 1 FROM chat_users
			WHERE chat_id = $1 AND user_id = $2
		);`,
		id, userID,
	)
	chat, err := scanRow[Chat](row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return chat, err
}

func (conn *MongoConnection) GetChatsByUserID(userID int64) ([]Chat, error) {
	return scanRows[Chat](conn.Query(
		`SELECT c.*
		FROM chat c
		INNER JOIN chat_users cu ON cu.chat_id = c.id
		LEFT JOIN (
			SELECT chat_id, MAX(created_at) AS last_message_time
			FROM message
			GROUP BY chat_id
		) m ON m.chat_id = c.id
		WHERE cu.user_id = $1
		ORDER BY COALESCE(m.last_message_time, c.last_updated_at) DESC;`,
		userID,
	))
}

// TODO: remove?
func (conn *MongoConnection) CheckPrivateChatExists(userIDs [2]int64) (bool, error) {
	result, err := conn.Exec(
		`SELECT 1
		FROM chat c
		JOIN chat_users cu1 ON cu1.chat_id = c.id AND cu1.user_id = $1
		JOIN chat_users cu2 ON cu2.chat_id = c.id AND cu2.user_id = $2
		WHERE c.type = $3`,
		userIDs[0], userIDs[1], PRIVATE_CHAT,
	)
	if err != nil {
		return false, err
	}
	chatCount, err := result.RowsAffected()
	return chatCount > 0, err
}

func (conn *MongoConnection) SetChat(chat *Chat, userIDs ...int64) (*Chat, error) {
	txn, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	chat, err = scanRow[Chat](txn.QueryRow(
		`INSERT INTO chat (type, name) VALUES ($1, $2)
		RETURNING *`,
		chat.Type, chat.Name,
	))
	if err != nil {
		return nil, errors.Join(err, txn.Rollback())
	}

	stmt, err := txn.Prepare(pq.CopyIn("chat_users", "chat_id", "user_id"))
	if err != nil {
		return nil, errors.Join(err, txn.Rollback())
	}

	for _, userID := range userIDs {
		_, err = stmt.Exec(chat.ID, userID)
		if err != nil {
			return nil, errors.Join(err, txn.Rollback())
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, errors.Join(err, txn.Rollback())
	}

	err = stmt.Close()
	if err != nil {
		return nil, errors.Join(err, txn.Rollback())
	}

	err = txn.Commit()
	if err != nil {
		return nil, err
	}
	return chat, nil
}
