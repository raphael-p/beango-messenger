package resolvers

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/raphael-p/beango/client"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func Home(w *response.Writer, r *http.Request, conn database.Connection) {
	user, _, httpError := getRequestContext(r)
	if ProcessHTTPError(w, httpError) {
		return
	}
	chats, httpError := chatsDatabase(user.ID, conn)
	if ProcessHTTPError(w, httpError) {
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

func OpenChat(w *response.Writer, r *http.Request, conn database.Connection) {
	user, params, httpError := getRequestContext(r, CHAT_ID_KEY)
	if ProcessHTTPError(w, httpError) {
		return
	}
	messages, httpError := chatMessagesDatabase(user.ID, params.ChatID, conn)
	if ProcessHTTPError(w, httpError) {
		return
	}
	chatName, err := context.GetParam(r, CHAT_NAME_KEY)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
	homeChat, err := template.New("home").Parse(client.MessagePane)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
	chatlist := map[string]any{"Title": chatName, "Messages": messages}
	homeChat.Execute(w, chatlist)
}
