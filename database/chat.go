package database

import (
	"database/sql"
	"time"
)

type chatType string

const (
	PRIVATE_CHAT chatType = "private"
	GROUP_CHAT   chatType = "group"
)

type Chat struct {
	ID            int64     `json:"id"`
	ChatType      chatType  `json:"chatType"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"LastUpdatedAt"`
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

func (conn *MongoConnection) CheckPrivateChatExists(userIDs [2]int64) bool {
	for _, chat := range Chats {
		if chat.ChatType == PRIVATE_CHAT {
			match := [2]bool{false, false}
			for _, chatUser := range ChatUsers {
				if chatUser.ChatID == chat.ID {
					if chatUser.UserID == userIDs[0] {
						if match[0] {
							break
						}
						match[0] = true
					} else if chatUser.UserID == userIDs[1] {
						if match[1] {
							break
						}
						match[1] = true
					} else {
						break
					}
				}
			}
			if match[0] && match[1] {
				return true
			}

		}
	}
	return false
}

func (conn *MongoConnection) SetChat(chat *Chat, userIDs ...int64) *Chat {
	chat.ID = int64(len(Chats) + 1)
	Chats[chat.ID] = *chat
	for _, userID := range userIDs {
		chatUser := ChatUser{
			ID:     int64(len(ChatUsers) + 1),
			ChatID: chat.ID,
			UserID: userID,
		}
		ChatUsers[chatUser.ID] = chatUser
	}
	return chat
}

// GPT:
// func (conn *MongoConnection) SetChat(chat *Chat, userIDs ...int) *Chat {
// 	// Insert chat into the 'chat' table
// 	_, err := conn.Exec(`
// 		INSERT INTO chat (type, name, created_at, last_updated_at)
// 		VALUES ($1, $2, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC')
// 		RETURNING id`,
// 		chat.Type, chat.Name)
// 	if err != nil {
// 		handleError(err)
// 		return nil
// 	}

// 	// Get the inserted chat ID
// 	var chatID int
// 	err = conn.QueryRow("SELECT lastval()").Scan(&chatID)
// 	if err != nil {
// 		handleError(err)
// 		return nil
// 	}
// 	chat.ID = chatID

// 	// Insert chat users into the 'chat_users' table
// 	for _, userID := range userIDs {
// 		_, err := conn.Exec(`
// 			INSERT INTO chat_users (chat_id, user_id, created_at)
// 			VALUES ($1, $2, NOW() AT TIME ZONE 'UTC')`,
// 			chat.ID, userID)
// 		if err != nil {
// 			handleError(err)
// 			return nil
// 		}
// 	}

// 	return chat
// }
