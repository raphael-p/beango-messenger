package client

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

// TODO: big refactor
// need to separate resolver into service + resolver
// move some client stuff to resolver
// e.g. MessageExtended becomes Message & Message becomes MessageDatabase

func Home(w *response.Writer, r *http.Request, conn database.Connection) {
	chats, ok := resolvers.GetChatsData(w, r, conn)
	if !ok {
		return
	}
	home, err := template.New("home").Parse(homePage)
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

	skeleton, err := getSkeleton()
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
	messages, ok := resolvers.GetChatMessagesData(w, r, conn)
	if !ok {
		return
	}
	chatName, err := context.GetParam(r, resolvers.CHAT_NAME_KEY)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
	homeChat, err := template.New("home").Parse(messagePane)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}
	chatlist := map[string]any{"Title": chatName, "Messages": messages}
	homeChat.Execute(w, chatlist)
}
