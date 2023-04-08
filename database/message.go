package database

type Message struct {
	ID      string `json:"id"`
	UserID  string `json:"userID"`
	ChatID  string `json:"chatID"`
	Content string `json:"content"`
}

func GetMessagesByChatID(chatID string) []Message {
	var messages []Message
	for _, message := range Messages {
		if message.ChatID == chatID {
			messages = append(messages, message)
		}
	}
	return messages
}

func SetMessage(message *Message) {
	Messages[message.ID] = *message
}
