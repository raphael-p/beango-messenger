package database

type Message struct {
	Id      string `json:"id"`
	UserId  string `json:"userId"`
	ChatId  string `json:"chatId"`
	Content string `json:"content"`
}

func GetMessagesByChatId(chatId string) []Message {
	var messages []Message
	for _, message := range Messages {
		if message.ChatId == chatId {
			messages = append(messages, message)
		}
	}
	return messages
}
