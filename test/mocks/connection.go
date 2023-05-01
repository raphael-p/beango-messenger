package mocks

import (
	"errors"
	"testing"

	"github.com/raphael-p/beango/database"
)

var Admin *database.User
var AdminSesh database.Session

var testUsers = make(map[string]database.User)
var testChats = make(map[string]database.Chat)
var testMessages = make(map[string]database.Message)
var testSessions = make(map[string]database.Session)

func populateMockDB(conn database.Connection) {
	Admin = MakeAdminUser()
	conn.SetUser(Admin)
	AdminSesh = MakeSession(Admin.ID)
	conn.SetSession(AdminSesh)
}

func clearMockDB() {
	testUsers = make(map[string]database.User)
	testChats = make(map[string]database.Chat)
	testMessages = make(map[string]database.Message)
	testSessions = make(map[string]database.Session)
}

type MockConnection struct{}

func MakeMockConnection(t *testing.T) database.Connection {
	conn := &MockConnection{}
	populateMockDB(conn)
	t.Cleanup(clearMockDB)
	return conn
}

func (tc *MockConnection) CheckSession(id string) (*database.Session, bool) {
	session := tc.GetSession(id)
	if session == nil {
		return nil, false
	} else {
		return session, true
	}
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

func (*MockConnection) GetSession(id string) *database.Session {
	session := testSessions[id]
	return &session
}

func (*MockConnection) GetSessionByUserID(userID string) (*database.Session, error) {
	return nil, nil
}

func (*MockConnection) GetUser(id string) (*database.User, error) {
	user, ok := testUsers[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &user, nil
}

func (*MockConnection) GetUserByUsername(username string) (*database.User, error) {
	return nil, nil
}

func (*MockConnection) SetChat(chat *database.Chat) {
	testChats[chat.ID] = *chat
}

func (*MockConnection) SetMessage(message *database.Message) {
	testMessages[message.ID] = *message
}

func (*MockConnection) SetSession(session database.Session) {
	testSessions[session.ID] = session
}

func (*MockConnection) SetUser(user *database.User) {
	testUsers[user.ID] = *user
}
