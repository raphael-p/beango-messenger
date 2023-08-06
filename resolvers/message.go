package resolvers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
)

// TODO: generic version of this
func ExtractUserAndChatIDFromRequest(w *response.Writer, r *http.Request) (int64, int64, bool) {
	user, params, httpError := getRequestContext(r, CHAT_ID_KEY)
	if ProcessHTTPError(w, httpError) {
		return 0, 0, false
	}
	chatID, httpErr := StringToInt(params[CHAT_ID_KEY], 64)
	if ProcessHTTPError(w, httpErr) {
		return 0, 0, false
	}
	return user.ID, chatID, true
}

// TODO: unit testing
func GetChatMessagesDatabase(userID, chatID int64, conn database.Connection) ([]database.MessageExtended, *HTTPError) {
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
	userID, chatID, ok := ExtractUserAndChatIDFromRequest(w, r)
	if !ok {
		return
	}
	messages, httpError := GetChatMessagesDatabase(userID, chatID, conn)
	fmt.Printf("messages: %v/n", messages)
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
		ProcessHTTPError(w, HandleDatabaseError(err))
		return
	}
	w.WriteJSON(http.StatusAccepted, newMessage)
}
