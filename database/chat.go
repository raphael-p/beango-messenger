package database

type Chat struct {
	Id      string    `json:"id"`
	UserIds [2]string `json:"userIds"`
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
