package resolvers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
)

func GetChatMessages(w *response.Writer, r *http.Request) {
	paramKeys := []string{"chatid"}
	user, params, ok := getRequestContext(w, r, true, paramKeys...)
	if !ok {
		return
	}
	chatId := params[paramKeys[0]]

	chat, _ := database.GetChat(chatId)
	if chat == nil || (chat.UserIds[0] != user.Id && chat.UserIds[1] != user.Id) {
		w.WriteString(http.StatusNotFound, "chat not found")
		return
	}

	w.WriteJSON(http.StatusOK, database.GetMessagesByChatId(chat.Id))
}

type SendMessageInput struct {
	Content string `json:"content"`
}

func SendMessage(w *response.Writer, r *http.Request) {
	paramKeys := []string{"chatid"}
	user, params, ok := getRequestContext(w, r, true, paramKeys...)
	if !ok {
		return
	}
	chatId := params[paramKeys[0]]

	chat, _ := database.GetChat(chatId)
	if chat == nil || (chat.UserIds[0] != user.Id && chat.UserIds[1] != user.Id) {
		w.WriteString(http.StatusNotFound, "chat not found")
		return
	}

	var input SendMessageInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	newMessage := &database.Message{
		Id:      uuid.NewString(),
		ChatId:  chatId,
		Content: input.Content,
	}
	database.SetMessage(newMessage)
	w.WriteJSON(http.StatusAccepted, newMessage)
}
