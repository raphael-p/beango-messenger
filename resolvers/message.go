package resolvers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
)

func GetChatMessages(w *response.Writer, r *http.Request) {
	params := struct{ chatID string }{}
	user, ok := getRequestContext(w, r, params)
	if !ok {
		return
	}

	chat, _ := database.GetChat(params.chatID)
	if chat == nil || (chat.UserIDs[0] != user.ID && chat.UserIDs[1] != user.ID) {
		w.WriteString(http.StatusNotFound, "chat not found")
		return
	}

	w.WriteJSON(http.StatusOK, database.GetMessagesByChatID(chat.ID))
}

type SendMessageInput struct {
	Content string `json:"content"`
}

func SendMessage(w *response.Writer, r *http.Request) {
	params := struct{ chatID string }{}
	user, ok := getRequestContext(w, r, params)
	if !ok {
		return
	}

	chat, _ := database.GetChat(params.chatID)
	if chat == nil || (chat.UserIDs[0] != user.ID && chat.UserIDs[1] != user.ID) {
		w.WriteString(http.StatusNotFound, "chat not found")
		return
	}

	var input SendMessageInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	newMessage := &database.Message{
		ID:      uuid.NewString(),
		ChatID:  params.chatID,
		Content: input.Content,
	}
	database.SetMessage(newMessage)
	w.WriteJSON(http.StatusAccepted, newMessage)
}
