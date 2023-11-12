package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/raphael-p/beango/config"
)

var Users = make(map[int64]User)
var Chats = make(map[int64]Chat)
var ChatUsers = make(map[int64]ChatUser)
var Messages = make(map[int64]MessageDatabase)
var Sessions = make(map[string]Session)

type Connection interface {
	GetChat(id, userID int64) (*Chat, error)
	GetChatsByUserID(userID int64) ([]Chat, error)
	GetPrivateChatByUserIDs(userID1, userID2 int64) (*Chat, error)
	SetChat(chat *Chat, userIDs ...int64) (*Chat, error)
	GetMessagesByChatID(chatID, fromMessageID, toMessageID int64, limit int) ([]Message, error)
	SetMessage(message *MessageDatabase) (*MessageDatabase, error)
	GetUser(id int64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUsersByChatID(chatID int64) ([]User, error)
	SetUser(user *User) (*User, error)
	SearchUsers(username string, searchUserID int64) ([]User, error)
	GetSession(id string) *Session
	GetSessionByUserID(userID int64) (*Session, error)
	SetSession(session Session)
	CheckSession(id string) (*Session, bool)
	DeleteSession(id string)
}

type MongoConnection struct {
	*sql.DB
}

var conn *MongoConnection

func GetConnection() (*MongoConnection, error) {
	if conn != nil {
		return conn, nil
	}

	// check for db envars
	host := os.Getenv(config.Envars.DatabaseHost)
	if host == "" {
		return nil, fmt.Errorf("$%s must be set", config.Envars.DatabaseHost)
	}
	name := os.Getenv(config.Envars.DatabaseName)
	if name == "" {
		return nil, fmt.Errorf("$%s must be set", config.Envars.DatabaseName)
	}

	// generate credentials substring
	credentials := ""
	username := os.Getenv(config.Envars.DatabaseUsername)
	password := os.Getenv(config.Envars.DatabasePassword)
	if username != "" && password != "" {
		credentials = username + ":" + password + "@"
	}

	connectionString := fmt.Sprintf("postgres://%s%s/%s?sslmode=disable", credentials, host, name)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	conn = &MongoConnection{db}
	return conn, nil
}

func SetDummyConnection() {
	conn = &MongoConnection{}
}
