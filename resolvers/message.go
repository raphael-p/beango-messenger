package resolvers

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/response"
)

func chatMessagesDatabase(userID, chatID int64, conn database.Connection) ([]database.Message, *resolverutils.HTTPError) {
	chat, err := conn.GetChat(chatID, userID)
	if err != nil {
		return nil, resolverutils.HandleDatabaseError(err)
	}
	if chat == nil {
		return nil, &resolverutils.HTTPError{http.StatusNotFound, "chat not found"}
	}

	messages, err := conn.GetMessagesByChatID(chatID)
	if err != nil {
		return nil, resolverutils.HandleDatabaseError(err)
	}
	return messages, nil
}

func GetChatMessages(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	messages, httpError := chatMessagesDatabase(user.ID, params.ChatID, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	w.WriteJSON(http.StatusOK, messages)
}

type SendMessageInput struct {
	Content string `json:"content"`
}

func SendMessage(w *response.Writer, r *http.Request, conn database.Connection) {
	var input SendMessageInput
	user, params, httpError := resolverutils.GetRequestBodyAndContext(r, &input, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	if chat, _ := conn.GetChat(params.ChatID, user.ID); chat == nil {
		w.WriteString(http.StatusNotFound, "chat not found")
		return
	}

	newMessage := &database.MessageDatabase{
		UserID:  user.ID,
		ChatID:  params.ChatID,
		Content: input.Content,
	}
	newMessage, err := conn.SetMessage(newMessage)
	if err != nil {
		resolverutils.ProcessHTTPError(w, resolverutils.HandleDatabaseError(err))
		return
	}
	w.WriteJSON(http.StatusAccepted, newMessage)
}
