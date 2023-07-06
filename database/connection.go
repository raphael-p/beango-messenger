package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/raphael-p/beango/config"
)

var Users = make(map[int64]User)
var Chats = make(map[int64]Chat)
var ChatUsers = make(map[int64]ChatUser)
var Messages = make(map[int64]Message)
var Sessions = make(map[string]Session)

type Connection interface {
	GetChat(id, userID int64) (*Chat, error)
	GetChatsByUserID(userID int64) []Chat
	CheckPrivateChatExists(userIDs [2]int64) bool
	SetChat(chat *Chat, userIDs ...int64) *Chat
	GetMessagesByChatID(chatID int64) []Message
	SetMessage(message *Message) (*Message, error)
	GetUser(id int64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	SetUser(user *User) (*User, error)
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

func SetDummyConnection() {
	conn = &MongoConnection{}
}
