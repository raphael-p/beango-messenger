package database

import (
	"fmt"
	"time"
)

type chatType string

const (
	PRIVATE_CHAT chatType = "private"
	GROUP_CHAT   chatType = "group"
)

type Chat struct {
	ID            int       `json:"id"`
	ChatType      chatType  `json:"chatType"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"LastUpdatedAt"`
}

type ChatUser struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chatID"`
	UserID    int       `json:"userID"`
	CreatedAt time.Time `json:"createdAt"`
}

func (conn *MongoConnection) GetChat(id, userID int) (*Chat, error) {
	chat, ok := Chats[id]
	if ok {
		for _, chatUser := range ChatUsers {
			if chatUser.UserID == userID && chatUser.ChatID == chat.ID {
				return &chat, nil
			}
		}
	}
	return nil, fmt.Errorf("no chat found with id %d", id)
}

func (conn *MongoConnection) GetChatsByUserID(userID int) []Chat {
	chats := []Chat{}
	for _, chatUser := range ChatUsers {
		if chatUser.UserID == userID {
			if chat, ok := Chats[chatUser.ChatID]; ok {
				chats = append(chats, chat)
			}
		}
	}
	return chats
}

func (conn *MongoConnection) CheckPrivateChatExists(userIDs [2]int) bool {
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

func (conn *MongoConnection) SetChat(chat *Chat, userIDs ...int) *Chat {
	chat.ID = len(Chats) + 1
	Chats[chat.ID] = *chat
	for _, userID := range userIDs {
		chatUser := ChatUser{
			ID:     len(ChatUsers) + 1,
			ChatID: chat.ID,
			UserID: userID,
		}
		ChatUsers[chatUser.ID] = chatUser
	}
	return chat
}
