package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/raphael-p/beango/config"
)

var Users = make(map[string]User)
var Chats = make(map[string]Chat)
var Messages = make(map[string]Message)
var Sessions = make(map[string]Session)

type Connection interface {
	CheckSession(id string) (*Session, bool)
	DeleteSession(id string)
	GetChat(id string) (*Chat, error)
	GetChatByUserIDs(userIDs [2]string) *Chat
	GetChatsByUserID(userID string) []Chat
	GetMessagesByChatID(chatID string) []Message
	GetSession(id string) *Session
	GetSessionByUserID(userID string) (*Session, error)
	GetUser(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	SetChat(chat *Chat)
	SetMessage(message *Message)
	SetSession(session Session)
	SetUser(user *User)
}

type MongoConnection struct {
	*sql.DB
}

func NewConnection() (*MongoConnection, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s/%s?sslmode=disable",
		config.Values.Database.Host,
		config.Values.Database.Port,
		config.Values.Database.Name,
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &MongoConnection{db}, nil
}
