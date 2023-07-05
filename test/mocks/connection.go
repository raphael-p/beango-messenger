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
	users     map[int64]database.User
	chats     map[int64]database.Chat
	chatUsers map[int64]database.ChatUser
	messages  map[int64]database.Message
	sessions  map[string]database.Session
}

func MakeMockConnection() *MockConnection {
	conn := &MockConnection{
		make(map[int64]database.User),
		make(map[int64]database.Chat),
		make(map[int64]database.ChatUser),
		make(map[int64]database.Message),
		make(map[string]database.Session),
	}
	populateMockDB(conn)
	return conn
}

func (mc *MockConnection) GetChat(id, userID int64) (*database.Chat, error) {
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

func (mc *MockConnection) GetChatsByUserID(userID int64) []database.Chat {
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

func (mc *MockConnection) CheckPrivateChatExists(userIDs [2]int64) bool {
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

func (mc *MockConnection) SetChat(chat *database.Chat, userIDs ...int64) *database.Chat {
	chat.ID = int64(len(mc.chats) + 1)
	mc.chats[chat.ID] = *chat
	for _, userID := range userIDs {
		chatUser := database.ChatUser{
			ID:     int64(len(mc.chatUsers) + 1),
			ChatID: chat.ID,
			UserID: userID,
		}
		mc.chatUsers[chatUser.ID] = chatUser
	}
	mc.chats[chat.ID] = *chat
	return chat
}

func (mc *MockConnection) GetMessagesByChatID(chatID int64) []database.Message {
	messages := []database.Message{}
	for _, message := range mc.messages {
		if message.ChatID == chatID {
			messages = append(messages, message)
		}
	}
	return messages
}

func (mc *MockConnection) SetMessage(message *database.Message) *database.Message {
	message.ID = int64(len(mc.messages) + 1)
	mc.messages[message.ID] = *message
	return message
}

func (mc *MockConnection) GetUser(id int64) (*database.User, error) {
	user, ok := mc.users[id]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

func (mc *MockConnection) GetUserByUsername(username string) (*database.User, error) {
	for _, user := range mc.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, nil
}

func (mc *MockConnection) SetUser(user *database.User) (*database.User, error) {
	user.ID = int64(len(mc.users) + 1)
	mc.users[user.ID] = *user
	return user, nil
}

func (mc *MockConnection) GetSession(id string) *database.Session {
	session, ok := mc.sessions[id]
	if !ok {
		return nil
	}
	return &session
}

func (*MockConnection) GetSessionByUserID(userID int64) (*database.Session, error) {
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
