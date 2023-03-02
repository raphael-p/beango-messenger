package resolvers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
)

func GetChats(w *utils.ResponseWriter, r *http.Request) {
	user := utils.GetUserFromContext(r)
	chats := database.GetChatsByUserId(user.Id)
	w.JSONResponse(http.StatusOK, chats)
}

type CreateChatInput struct {
	UserId string `json:"userId"`
}

func CreateChat(w *utils.ResponseWriter, r *http.Request) {
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

	user := utils.GetUserFromContext(r)
	userIds := [2]string{user.Id, input.UserId}

	// Check if chat already exists
	if chat := database.GetChatByUserIds(userIds); chat != nil {
		w.StringResponse(http.StatusConflict, "chat already exists")
		return
	}

	newChat := database.Chat{
		Id:      uuid.NewString(),
		UserIds: userIds,
	}
	w.JSONResponse(http.StatusCreated, newChat)
}
