package mocks

import (
	"errors"

	"github.com/raphael-p/beango/database"
)

var Admin *database.User
var AdminSesh database.Session

func populateMockDB(conn database.Connection) {
	Admin = MakeAdminUser()
	conn.SetUser(Admin)
	AdminSesh = MakeSession(Admin.ID)
	conn.SetSession(AdminSesh)
}

type MockConnection struct {
	users    map[string]database.User
	chats    map[string]database.Chat
	messages map[string]database.Message
	sessions map[string]database.Session
}

func MakeMockConnection() database.Connection {
	conn := &MockConnection{
		make(map[string]database.User),
		make(map[string]database.Chat),
		make(map[string]database.Message),
		make(map[string]database.Session),
	}
	populateMockDB(conn)
	return conn
}

func (mc *MockConnection) CheckSession(id string) (*database.Session, bool) {
	session := mc.GetSession(id)
	if session == nil {
		return nil, false
	}
	return session, true
}

func (*MockConnection) DeleteSession(id string) {}

func (*MockConnection) GetChat(id string) (*database.Chat, error) {
	return nil, nil
}

func (*MockConnection) GetChatByUserIDs(userIDs [2]string) *database.Chat {
	return nil
}

func (*MockConnection) GetChatsByUserID(userID string) []database.Chat {
	return nil
}

func (*MockConnection) GetMessagesByChatID(chatID string) []database.Message {
	return nil
}

func (mc *MockConnection) GetSession(id string) *database.Session {
	session, ok := mc.sessions[id]
	if !ok {
		return nil
	}
	return &session
}

func (*MockConnection) GetSessionByUserID(userID string) (*database.Session, error) {
	return nil, nil
}

func (mc *MockConnection) GetUser(id string) (*database.User, error) {
	user, ok := mc.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &user, nil
}

func (mc *MockConnection) GetUserByUsername(username string) (*database.User, error) {
	for _, user := range mc.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("not found")
}

func (mc *MockConnection) SetChat(chat *database.Chat) {
	mc.chats[chat.ID] = *chat
}

func (mc *MockConnection) SetMessage(message *database.Message) {
	mc.messages[message.ID] = *message
}

func (mc *MockConnection) SetSession(session database.Session) {
	mc.sessions[session.ID] = session
}

func (mc *MockConnection) SetUser(user *database.User) {
	mc.users[user.ID] = *user
}
