package database

import (
	"fmt"
)

type Chat struct {
	Id      string    `json:"id"`
	UserIds [2]string `json:"userIds"`
}

func GetChat(id string) (*Chat, error) {
	chat, ok := Chats[id]
	if !ok {
		return nil, fmt.Errorf("no chat found with id %s", id)
	} else {
		return &chat, nil
	}
}

func GetChatsByUserId(userId string) []Chat {
	var chats []Chat
	for _, chat := range Chats {
		for _, chatUserId := range chat.UserIds {
			if chatUserId == userId {
				chats = append(chats, chat)
			}
		}
	}
	return chats
}

func GetChatByUserIds(userIds [2]string) *Chat {
	for _, chat := range Chats {
		if (chat.UserIds[0] == userIds[0] &&
			chat.UserIds[1] == userIds[1]) ||
			(chat.UserIds[0] == userIds[1] &&
				chat.UserIds[1] == userIds[0]) {
			return &chat
		}

	}
	return nil
}
