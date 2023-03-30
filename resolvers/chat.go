package resolvers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/httputils"
)

func GetChats(w *httputils.ResponseWriter, r *http.Request) {
	user, err := httputils.GetContextUser(r)
	if err != nil {
		w.StringResponse(http.StatusInternalServerError, err.Error())
		return
	}
	chats := database.GetChatsByUserId(user.Id)
	w.JSONResponse(http.StatusOK, chats)
}

type CreateChatInput struct {
	UserId string `json:"userId"`
}

func CreateChat(w *httputils.ResponseWriter, r *http.Request) {
	var input CreateChatInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	// Check that user id exists
	_, err := database.GetUser(input.UserId)
	if err != nil {
		message := fmt.Sprintf("userId %s is invalid", input.UserId)
		w.StringResponse(http.StatusBadRequest, message)
		return
	}

	user, err := httputils.GetContextUser(r)
	if err != nil {
		w.StringResponse(http.StatusInternalServerError, err.Error())
		return
	}
	userIds := [2]string{user.Id, input.UserId}

	// Check if chat already exists
	if chat := database.GetChatByUserIds(userIds); chat != nil {
		w.StringResponse(http.StatusConflict, "chat already exists")
		return
	}

	newChat := &database.Chat{
		Id:      uuid.NewString(),
		UserIds: userIds,
	}
	database.SetChat(newChat)
	w.JSONResponse(http.StatusCreated, newChat)
}
