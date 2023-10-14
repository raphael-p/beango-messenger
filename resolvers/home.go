package resolvers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/raphael-p/beango/client"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/response"
)

func Home(w *response.Writer, r *http.Request, conn database.Connection) {
	user, _, httpError := resolverutils.GetRequestContext(r)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	chats, httpError := chatsDatabase(user.ID, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	chatlist := map[string]any{"Chats": chats}
	client.ServeTemplate(w, "homePage", client.Skeleton+client.HomePage, chatlist)
}

// TODO: paging for messages
//
//	 -> check that "from" param excludes last
//		-> add "to" param. should exclude that value
//		-> add "limit" param
//
// TODO: unit testing for home.go + login.go
func OpenChat(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	messages, lastMessageID, httpError := getMessages(user.ID, params.ChatID, 0, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	chatName, httpError := resolverutils.GetRequestQueryParam(r, "name", true)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	chatlist := map[string]any{"Name": chatName, "Messages": messages, "ID": params.ChatID, "FromMessageID": lastMessageID}
	client.ServeTemplate(w, "messagePane", client.MessagePane, chatlist)
}

func RefreshChat(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	fromMessageIDParam, httpError := resolverutils.GetRequestQueryParam(r, "from", true)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	var fromMessageID int64
	if fromMessageIDParam != "" {
		var err error
		fromMessageID, err = strconv.ParseInt(fromMessageIDParam, 10, 64)
		if err != nil {
			w.WriteString(http.StatusBadRequest, "query parameter 'from' must be an integer")
			return
		}
	}

	messages, lastMessageID, httpError := getMessages(user.ID, params.ChatID, fromMessageID, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	if len(messages) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	chatlist := map[string]any{"Messages": messages, "ID": params.ChatID, "FromMessageID": lastMessageID}
	client.ServeTemplate(w, "messagePaneRefresh", client.MessagePaneRefresh, chatlist)
}

func getMessages(userID, chatID, fromMessageID int64, conn database.Connection) ([]database.Message, int64, *resolverutils.HTTPError) {
	messages, httpError := chatMessagesDatabase(userID, chatID, fromMessageID, conn)
	if httpError != nil {
		return nil, 0, httpError
	}

	var lastMessageID int64
	if len(messages) != 0 {
		lastMessageIndex := int(math.Max(float64(len(messages)-1), 0))
		lastMessageID = messages[lastMessageIndex].ID
	}
	return messages, lastMessageID, nil
}

func SendChatMessage(w *response.Writer, r *http.Request, conn database.Connection) {
	content := r.PostFormValue("content")
	if content == "" {
		w.WriteString(http.StatusBadRequest, "cannot send an empty message")
		return
	}
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	_, httpError = sendMessageDatabase(user.ID, params.ChatID, content, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	w.Header().Set("HX-Trigger", "chat-refresh")
	w.WriteHeader(http.StatusNoContent)
}
