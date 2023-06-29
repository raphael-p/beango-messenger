package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/raphael-p/beango/config"
)

var Users = make(map[int]User)
var Chats = make(map[int]Chat)
var ChatUsers = make(map[int]ChatUser)
var Messages = make(map[int]Message)
var Sessions = make(map[string]Session)

type Connection interface {
	GetChat(id, userID int) (*Chat, error)
	GetChatsByUserID(userID int) []Chat
	CheckPrivateChatExists(userIDs [2]int) bool
	SetChat(chat *Chat, userIDs ...int) *Chat
	GetMessagesByChatID(chatID int) []Message
	SetMessage(message *Message) *Message
	GetUser(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
	SetUser(user *User) *User
	GetSession(id string) *Session
	GetSessionByUserID(userID int) (*Session, error)
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
