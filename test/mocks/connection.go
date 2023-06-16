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
	users     map[int]database.User
	chats     map[int]database.Chat
	chatUsers map[int]database.ChatUser
	messages  map[int]database.Message
	sessions  map[string]database.Session
}

func MakeMockConnection() *MockConnection {
	conn := &MockConnection{
		make(map[int]database.User),
		make(map[int]database.Chat),
		make(map[int]database.ChatUser),
		make(map[int]database.Message),
		make(map[string]database.Session),
	}
	populateMockDB(conn)
	return conn
}

func (mc *MockConnection) GetChat(id, userID int) (*database.Chat, error) {
	chat, ok := mc.chats[id]
	if ok {
		for _, chatUser := range mc.chatUsers {
			if chatUser.UserID == userID && chatUser.ChatID == chat.ID {
				return &chat, nil
			}
		}
	}
	return nil, errors.New("not found")
}

func (mc *MockConnection) GetChatsByUserID(userID int) []database.Chat {
	chats := []database.Chat{}
	for _, chatUser := range mc.chatUsers {
		if chatUser.UserID == userID {
			if chat, ok := mc.chats[chatUser.ChatID]; ok {
				chats = append(chats, chat)
			}
		}
	}
	return chats
}

func (mc *MockConnection) CheckPrivateChatExists(userIDs [2]int) bool {
	for _, chat := range mc.chats {
		if chat.ChatType == database.PRIVATE_CHAT {
			match := [2]bool{false, false}
			for _, chatUser := range mc.chatUsers {
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

func (mc *MockConnection) SetChat(chat *database.Chat, userIDs ...int) {
	chat.ID = len(mc.chats) + 1
	mc.chats[chat.ID] = *chat
	for _, userID := range userIDs {
		chatUser := database.ChatUser{
			ID:     len(mc.chatUsers) + 1,
			ChatID: chat.ID,
			UserID: userID,
		}
		mc.chatUsers[chatUser.ID] = chatUser
	}
	mc.chats[chat.ID] = *chat
}

func (mc *MockConnection) GetMessagesByChatID(chatID int) []database.Message {
	messages := []database.Message{}
	for _, message := range mc.messages {
		if message.ChatID == chatID {
			messages = append(messages, message)
		}
	}
	return messages
}

func (mc *MockConnection) SetMessage(message *database.Message) {
	mc.messages[message.ID] = *message
}

func (mc *MockConnection) GetUser(id int) (*database.User, error) {
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

func (mc *MockConnection) SetUser(user *database.User) {
	mc.users[user.ID] = *user
}

func (mc *MockConnection) GetSession(id string) *database.Session {
	session, ok := mc.sessions[id]
	if !ok {
		return nil
	}
	return &session
}

func (*MockConnection) GetSessionByUserID(userID int) (*database.Session, error) {
	return nil, nil
}

func (mc *MockConnection) SetSession(session database.Session) {
	mc.sessions[session.ID] = session
}

func (mc *MockConnection) CheckSession(id string) (*database.Session, bool) {
	session := mc.GetSession(id)
	if session == nil {
		return nil, false
	}
	return session, true
}

func (mc *MockConnection) DeleteSession(id string) {
	delete(mc.sessions, id)
}
