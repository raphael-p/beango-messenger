package database

import (
	"fmt"
	"time"
)

type Chat struct {
	ID            string    `json:"id"`
	UserIDs       [2]string `json:"userIDs"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"LastUpdatedAt"`
}

func (conn *MongoConnection) GetChat(id string) (*Chat, error) {
	chat, ok := Chats[id]
	if !ok {
		return nil, fmt.Errorf("no chat found with id %s", id)
	} else {
		return &chat, nil
	}
}

func (conn *MongoConnection) GetChatsByUserID(userID string) []Chat {
	chats := []Chat{}
	for _, chat := range Chats {
		for _, chatUserID := range chat.UserIDs {
			if chatUserID == userID {
				chats = append(chats, chat)
			}
		}
	}
	return chats
}

func (conn *MongoConnection) GetChatByUserIDs(userIDs [2]string) *Chat {
	for _, chat := range Chats {
		if (chat.UserIDs[0] == userIDs[0] &&
			chat.UserIDs[1] == userIDs[1]) ||
			(chat.UserIDs[0] == userIDs[1] &&
				chat.UserIDs[1] == userIDs[0]) {
			return &chat
		}

	}
	return nil
}

func (conn *MongoConnection) SetChat(chat *Chat) {
	Chats[chat.ID] = *chat
}
