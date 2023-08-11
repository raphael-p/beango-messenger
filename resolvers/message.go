package resolvers

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
)

// TODO: unit testing
func chatMessagesDatabase(userID, chatID int64, conn database.Connection) ([]database.MessageExtended, *HTTPError) {
	chat, err := conn.GetChat(chatID, userID)
	if err != nil {
		return nil, HandleDatabaseError(err)
	}
	if chat == nil {
		return nil, &HTTPError{http.StatusNotFound, "chat not found"}
	}

	messages, err := conn.GetMessagesByChatID(chatID)
	if err != nil {
		return nil, HandleDatabaseError(err)
	}
	return messages, nil
}

func GetChatMessages(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := getRequestContext(r, CHAT_ID_KEY)
	if ProcessHTTPError(w, httpError) {
		return
	}
	messages, httpError := chatMessagesDatabase(user.ID, params.ChatID, conn)
	if ProcessHTTPError(w, httpError) {
		return
	}
	w.WriteJSON(http.StatusOK, messages)
}

type SendMessageInput struct {
	Content string `json:"content"`
}

func SendMessage(w *response.Writer, r *http.Request, conn database.Connection) {
	var input SendMessageInput
	user, params, httpError := getRequestBodyAndContext(r, &input, CHAT_ID_KEY)
	if ProcessHTTPError(w, httpError) {
		return
	}

	if chat, _ := conn.GetChat(params.ChatID, user.ID); chat == nil {
		w.WriteString(http.StatusNotFound, "chat not found")
		return
	}

	newMessage := &database.Message{
		UserID:  user.ID,
		ChatID:  params.ChatID,
		Content: input.Content,
	}
	newMessage, err := conn.SetMessage(newMessage)
	if err != nil {
		ProcessHTTPError(w, HandleDatabaseError(err))
		return
	}
	w.WriteJSON(http.StatusAccepted, newMessage)
}
