package resolvers

import (
	"net/http"
	"strings"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/response"
)

func chatMessagesDatabase(userID, chatID, fromMessageID, toMessageID int64, limit int, conn database.Connection) ([]database.Message, *resolverutils.HTTPError) {
	chat, err := conn.GetChat(chatID, userID)
	if err != nil {
		return nil, resolverutils.HandleDatabaseError(err)
	}
	if chat == nil {
		return nil, &resolverutils.HTTPError{
			Status:  http.StatusNotFound,
			Message: "chat not found",
		}
	}

	messages, err := conn.GetMessagesByChatID(chatID, fromMessageID, toMessageID, limit)
	return messages, resolverutils.HandleDatabaseError(err)
}

func GetChatMessages(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	messages, httpError := chatMessagesDatabase(user.ID, params.ChatID, 0, 0, 0, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	w.WriteJSON(http.StatusOK, messages)
}

type sendMessageInput struct {
	Content string `json:"content"`
}

func sendMessageDatabase(userID, chatID int64, content string, conn database.Connection) (*database.MessageDatabase, *resolverutils.HTTPError) {
	if chat, _ := conn.GetChat(chatID, userID); chat == nil {
		return nil, &resolverutils.HTTPError{
			Status:  http.StatusNotFound,
			Message: "chat not found",
		}
	}

	newMessage := &database.MessageDatabase{
		UserID:  userID,
		ChatID:  chatID,
		Content: strings.TrimSpace(content),
	}
	newMessage, err := conn.SetMessage(newMessage)
	return newMessage, resolverutils.HandleDatabaseError(err)
}

func SendMessage(w *response.Writer, r *http.Request, conn database.Connection) {
	var input sendMessageInput
	user, params, httpError := resolverutils.GetRequestBodyAndContext(r, &input, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	newMessage, httpError := sendMessageDatabase(user.ID, params.ChatID, input.Content, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	w.WriteJSON(http.StatusAccepted, newMessage)
}
