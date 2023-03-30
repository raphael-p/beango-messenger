package resolvers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/httputils"
)

func GetChatMessages(w *httputils.ResponseWriter, r *http.Request) {
	user, err := httputils.GetContextUser(r)
	if err != nil {
		w.StringResponse(http.StatusInternalServerError, err.Error())
		return
	}
	chatId, err := httputils.GetContextParam(r, "chatid")
	if err != nil {
		w.StringResponse(http.StatusInternalServerError, err.Error())
		return
	}
	chat, _ := database.GetChat(chatId)
	if chat == nil || (chat.UserIds[0] != user.Id && chat.UserIds[1] != user.Id) {
		w.StringResponse(http.StatusNotFound, "chat not found")
		return
	}

	w.JSONResponse(http.StatusOK, database.GetMessagesByChatId(chat.Id))
}

type SendMessageInput struct {
	Content string `json:"content"`
}

func SendMessage(w *httputils.ResponseWriter, r *http.Request) {
	user, err := httputils.GetContextUser(r)
	if err != nil {
		w.StringResponse(http.StatusInternalServerError, err.Error())
		return
	}
	chatId, err := httputils.GetContextParam(r, "chatid")
	if err != nil {
		w.StringResponse(http.StatusInternalServerError, err.Error())
		return
	}
	chat, _ := database.GetChat(chatId)
	if chat == nil || (chat.UserIds[0] != user.Id && chat.UserIds[1] != user.Id) {
		w.StringResponse(http.StatusNotFound, "chat not found")
		return
	}

	var input SendMessageInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	newMessage := &database.Message{
		Id:      uuid.NewString(),
		ChatId:  chatId,
		Content: input.Content,
	}
	database.SetMessage(newMessage)
	w.JSONResponse(http.StatusAccepted, newMessage)
}
