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

	chatList := map[string]any{"Chats": chats}
	client.ServeTemplate(w, "homePage", client.Skeleton+client.HomePage, chatList)
}

func OpenChat(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	chatName, httpError := resolverutils.GetRequestQueryParam(r, "name", true)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	chatData, httpError := openChatData(user.ID, params.ChatID, chatName, conn)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}
	client.ServeTemplate(w, "messagePane", client.MessagePane, chatData)
}

func openChatData(userID, chatID int64, chatName string, conn database.Connection) (map[string]any, *resolverutils.HTTPError) {
	messages, firstMessageID, lastMessageID, httpError := getMessages(
		userID,
		chatID,
		0,
		0,
		MESSAGE_BATCH_SIZE,
		conn,
	)
	if httpError != nil {
		return nil, httpError
	}

	data := map[string]any{
		"Name":          chatName,
		"Messages":      messages,
		"ID":            chatID,
		"FromMessageID": lastMessageID,
		"ToMessageID":   firstMessageID,
	}
	return data, nil
}

func RefreshMessages(w *response.Writer, r *http.Request, conn database.Connection) {
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

	chatData := map[string]any{
		"Messages":      messages,
		"ID":            params.ChatID,
		"FromMessageID": lastMessageID,
		"ToMessageID":   firstMessageID,
		"IsRefresh":     true,
	}
	client.ServeTemplate(w, "messagePaneRefresh", client.MessagePaneRefresh, chatData)
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

	olderMessages := map[string]any{
		"Messages":    messages,
		"ID":          params.ChatID,
		"ToMessageID": firstMessageID,
	}
	client.ServeTemplate(w, "messagePaneScroll", client.MessagePaneScroll, olderMessages)
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

type sendMessageHTMLInput struct {
	Content validate.JSONField[string] `json:"content" zeroable:"true"`
}

func SendMessageHTML(w *response.Writer, r *http.Request, conn database.Connection) {
	input := new(sendMessageHTMLInput)
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
	w.Header().Set("HX-Trigger", "refresh-messages")
	w.WriteHeader(http.StatusNoContent)
}

func OpenChatCreator(w *response.Writer, r *http.Request, conn database.Connection) {
	w.WriteString(http.StatusOK, client.NewChatPane)
}

type userSearchInput struct {
	Query validate.JSONField[string] `json:"query" zeroable:"true"`
}

func UserSearch(w *response.Writer, r *http.Request, conn database.Connection) {
	input := new(userSearchInput)
	user, _, httpError := resolverutils.GetRequestBodyAndContext(r, input)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	if input.Query.Value == "" {
		w.WriteString(http.StatusOK, "")
		return
	}

	users, err := conn.SearchUsers(input.Query.Value, user.ID)
	if resolverutils.DisplayHTTPError(w, resolverutils.HandleDatabaseError(err)) {
		return
	}

	data := map[string]any{"Users": stripUserFields(users...)}
	client.ServeTemplate(w, "userSearchResults", client.UserSearchResults, data)
}

func CreatePrivateChatHTML(w *response.Writer, r *http.Request, conn database.Connection) {
	var input createPrivateChatInput
	user, _, httpError := resolverutils.GetRequestBodyAndContext(r, &input)
	if resolverutils.DisplayHTTPError(w, httpError) ||
		resolverutils.DisplayHTTPError(w, validateCreatePrivateChatInput(&input, user.ID)) {
		return
	}

	newChat, httpError := createPrivateChatDatabase(user.ID, input.UserID, conn)
	if newChat == nil && resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	// Resolve chat display name
	chatName := newChat.Name
	if chatName == "" {
		inputUser, err := conn.GetUser(input.UserID)
		if resolverutils.DisplayHTTPError(w, resolverutils.HandleDatabaseError(err)) {
			return
		}
		chatName = generateChatName(user.ID, []database.User{*user, *inputUser})
	}

	chatData, httpError := openChatData(user.ID, newChat.ID, chatName, conn)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	// Get chat list with new chat
	chats, httpError := getChatsDatabase(user.ID, conn)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}
	chatData["Chats"] = chats

	client.ServeTemplate(w, "messagePaneWithChatRefresh", client.MessagePane+client.ChatListRefresh, chatData)
}
