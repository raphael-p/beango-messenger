package resolvers

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"

	"github.com/raphael-p/beango/client"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/logger"
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
	home, err := template.New("home").Parse(client.HomePage)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
	homeWithChatlist := new(bytes.Buffer)
	chatlist := map[string]any{"Chats": chats}
	home.Execute(homeWithChatlist, chatlist)

	if r.Header.Get("HX-Request") == "true" {
		w.Write([]byte("<div id='content' hx-swap-oob='innerHTML'>" + homeWithChatlist.String() + "</div>"))
		return
	}

	skeleton, err := client.GetSkeleton()
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
	data := map[string]any{"content": template.HTML(homeWithChatlist.String())}
	if err := skeleton.Execute(w, data); err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
	}
}

// TODO: only fetch new messages + pagination
// TODO: regular chat refresh without clearing message bar
// TODO: investigate adding message directly after sending (instead of triggering update event)
// TODO: refactor into 2 separate endpoints
// TODO: use hx-swap to automatically scroll to the bottom of messages on opening a chat (but not on refresh)
func OpenChat(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	isRefresh := r.URL.Query().Has("refresh")
	fromMessageIDParam, httpError := resolverutils.GetRequestQueryParam(r, "from", isRefresh)
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

	messages, httpError := chatMessagesDatabase(user.ID, params.ChatID, fromMessageID, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	lastMessageID := messages[len(messages)-1].ID

	chatName, httpError := resolverutils.GetRequestQueryParam(r, "name", !isRefresh)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	var templateToParse string
	if isRefresh {
		templateToParse = client.MessagePaneRefresh
	} else {
		templateToParse = client.MessagePane
	}
	chatTemplate, err := template.New("home").Parse(templateToParse)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
	chatlist := map[string]any{"Name": chatName, "Messages": messages, "ID": params.ChatID, "FromMessageID": lastMessageID}
	if err := chatTemplate.Execute(w, chatlist); err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
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
