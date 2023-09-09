package resolvers

import (
	"bytes"
	"html/template"
	"net/http"

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
func OpenChat(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	messages, httpError := chatMessagesDatabase(user.ID, params.ChatID, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	chatName, httpError := resolverutils.GetRequestQueryParam(r, "name", true, true)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	homeChat, err := template.New("home").Parse(client.MessagePane)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("HX-Trigger-After-Settle", "chat-opened") // has to be set before .Execute()
	chatlist := map[string]any{"Name": chatName, "Messages": messages, "ID": params.ChatID}
	if err := homeChat.Execute(w, chatlist); err != nil {
		w.Header().Del("HX-Trigger-After-Settle")
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
}
