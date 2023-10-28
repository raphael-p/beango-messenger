package resolvers

import (
	"net/http"
	"slices"

	"github.com/raphael-p/beango/client"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/response"
	"github.com/raphael-p/beango/utils/validate"
)

const MESSAGE_BATCH_SIZE int = 50

func Home(w *response.Writer, r *http.Request, conn database.Connection) {
	user, _, httpError := resolverutils.GetRequestContext(r)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	chats, httpError := getChatsDatabase(user.ID, conn)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	chatlist := map[string]any{"Chats": chats}
	client.ServeTemplate(w, "homePage", client.Skeleton+client.HomePage, chatlist)
}

func OpenChat(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	messages, firstMessageID, lastMessageID, httpError := getMessages(
		user.ID,
		params.ChatID,
		0,
		0,
		MESSAGE_BATCH_SIZE,
		conn,
	)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	chatName, httpError := resolverutils.GetRequestQueryParam(r, "name", true)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	chatlist := map[string]any{
		"Name":          chatName,
		"Messages":      messages,
		"ID":            params.ChatID,
		"FromMessageID": lastMessageID,
		"ToMessageID":   firstMessageID,
	}
	client.ServeTemplate(w, "messagePane", client.MessagePane, chatlist)
}

func RefreshChat(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	fromMessageID, httpError := resolverutils.GetRequestQueryParamInt(r, "from", true)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	messages, firstMessageID, lastMessageID, httpError := getMessages(
		user.ID,
		params.ChatID,
		fromMessageID,
		0,
		0,
		conn,
	)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}
	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	chatlist := map[string]any{
		"Messages":      messages,
		"ID":            params.ChatID,
		"FromMessageID": lastMessageID,
		"ToMessageID":   firstMessageID,
		"IsRefresh":     true,
	}
	client.ServeTemplate(w, "messagePaneRefresh", client.MessagePaneRefresh, chatlist)
}

func ScrollUp(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	toMessageID, httpError := resolverutils.GetRequestQueryParamInt(r, "to", true)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	messages, firstMessageID, _, httpError := getMessages(
		user.ID,
		params.ChatID,
		0,
		toMessageID,
		MESSAGE_BATCH_SIZE,
		conn,
	)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	chatlist := map[string]any{
		"Messages":    messages,
		"ID":          params.ChatID,
		"ToMessageID": firstMessageID,
	}
	client.ServeTemplate(w, "messagePaneScroll", client.MessagePaneScroll, chatlist)
}

func getMessages(userID, chatID, fromMessageID, toMessageID int64, limit int, conn database.Connection) ([]database.Message, int64, int64, *resolverutils.HTTPError) {
	messages, httpError := chatMessagesDatabase(userID, chatID, fromMessageID, toMessageID, limit, conn)
	if httpError != nil {
		return nil, 0, 0, httpError
	}

	var lastMessageID int64
	var firstMessageID int64
	if len(messages) != 0 {
		lastMessageID = messages[0].ID
		slices.Reverse(messages) // comes in sorted newest to oldest
		firstMessageID = messages[0].ID
	}

	return messages, firstMessageID, lastMessageID, nil
}

type sendChatMessageInput struct {
	Content validate.JSONField[string] `json:"content" zeroable:"true"`
}

func SendChatMessage(w *response.Writer, r *http.Request, conn database.Connection) {
	input := new(sendChatMessageInput)
	user, params, httpError := resolverutils.GetRequestBodyAndContext(r, input, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	if input.Content.Value == "" {
		w.WriteString(http.StatusBadRequest, "cannot send an empty message")
		return
	}

	_, httpError = sendMessageDatabase(user.ID, params.ChatID, input.Content.Value, conn)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}
	w.Header().Set("HX-Trigger", "chat-refresh")
	w.WriteHeader(http.StatusNoContent)
}
