package resolvers

import (
	"net/http"
	"strconv"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
)

func GetChatMessages(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, ok := getRequestContext(w, r, CHAT_ID_KEY)
	if !ok {
		return
	}

	chatID, err := strconv.Atoi(params[CHAT_ID_KEY])
	if err != nil {
		w.WriteString(http.StatusBadRequest, "chat ID must be an integer")
	}

	chat, _ := conn.GetChat(chatID, user.ID)
	if chat == nil {
		w.WriteString(http.StatusNotFound, "chat not found")
		return
	}

	w.WriteJSON(http.StatusOK, conn.GetMessagesByChatID(chat.ID))
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

	chatID, err := strconv.Atoi(params[CHAT_ID_KEY])
	if err != nil {
		w.WriteString(http.StatusBadRequest, "chat ID must be an integer")
	}

	chat, _ := conn.GetChat(chatID, user.ID)
	if chat == nil {
		w.WriteString(http.StatusNotFound, "chat not found")
		return
	}

	newMessage := &database.Message{
		UserID:  user.ID,
		ChatID:  chatID,
		Content: input.Content,
	}
	newMessage = conn.SetMessage(newMessage)
	w.WriteJSON(http.StatusAccepted, newMessage)
}
