package resolvers

import (
	"net/http"
	"strconv"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
)

func GetChatMessagesData(w *response.Writer, r *http.Request, conn database.Connection) ([]database.MessageExtended, bool) {
	user, params, ok := getRequestContext(w, r, CHAT_ID_KEY)
	if !ok {
		return nil, false
	}

	chatID, err := strconv.ParseInt(params[CHAT_ID_KEY], 10, 64)
	if err != nil {
		w.WriteString(http.StatusBadRequest, "chat ID must be an integer")
	}

	chat, err := conn.GetChat(chatID, user.ID)
	if err != nil {
		HandleDatabaseError(w, err)
		return nil, false
	}
	if chat == nil {
		w.WriteString(http.StatusNotFound, "chat not found")
		return nil, false
	}

	messages, err := conn.GetMessagesByChatID(chatID)
	if err != nil {
		HandleDatabaseError(w, err)
		return nil, false
	}
	return messages, true
}

func GetChatMessages(w *response.Writer, r *http.Request, conn database.Connection) {
	if messages, ok := GetChatMessagesData(w, r, conn); ok {
		w.WriteJSON(http.StatusOK, messages)
	}
}

type SendMessageInput struct {
	Content string `json:"content"`
}

func SendMessage(w *response.Writer, r *http.Request, conn database.Connection) {
	var input SendMessageInput
	user, params, ok := getRequestBodyAndContext(w, r, &input, CHAT_ID_KEY)
	if !ok {
		return
	}

	chatID, err := strconv.ParseInt(params[CHAT_ID_KEY], 10, 64)
	if err != nil {
		w.WriteString(http.StatusBadRequest, "chat ID must be an integer")
	}

	if chat, _ := conn.GetChat(chatID, user.ID); chat == nil {
		w.WriteString(http.StatusNotFound, "chat not found")
		return
	}

	newMessage := &database.Message{
		UserID:  user.ID,
		ChatID:  chatID,
		Content: input.Content,
	}
	newMessage, err = conn.SetMessage(newMessage)
	if err != nil {
		HandleDatabaseError(w, err)
		return
	}
	w.WriteJSON(http.StatusAccepted, newMessage)
}
