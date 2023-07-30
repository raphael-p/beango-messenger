package database

import (
	"database/sql"
	"time"
)

type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	DisplayName   string    `json:"displayName"`
	Key           []byte    `json:"key"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"LastUpdatedAt"`
}

func (conn *MongoConnection) GetUser(id int64) (*User, error) {
	row := conn.QueryRow(
		`SELECT * FROM "user" WHERE id = $1`,
		id,
	)
	user, err := scanRow[User](row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (conn *MongoConnection) GetUserByUsername(username string) (*User, error) {
	row := conn.QueryRow(
		`SELECT * FROM "user" WHERE username = $1
		ORDER BY created_at ASC 
		LIMIT 1`,
		username,
	)
	user, err := scanRow[User](row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (conn *MongoConnection) GetUsersByChatID(chatID int64) ([]User, error) {
	return scanRows[User](conn.Query(
		`SELECT u.*
		FROM "user" u
		INNER JOIN chat_users cu ON cu.user_id = u.id
		WHERE cu.chat_id = $1`,
		chatID,
	))
}

func (conn *MongoConnection) SetUser(user *User) (*User, error) {
	return scanRow[User](conn.QueryRow(
		`INSERT INTO "user" (username, display_name, key)
		VALUES ($1, $2, $3)
		RETURNING *`,
		user.Username, user.DisplayName, user.Key,
	))
}
