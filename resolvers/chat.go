package resolvers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
)

func GetChats(w *response.Writer, r *http.Request, conn database.Connection) {
	user, ok := getRequestContext(w, r, &struct{}{})
	if !ok {
		return
	}
	chats := conn.GetChatsByUserID(user.ID)
	w.WriteJSON(http.StatusOK, chats)
}

type CreateChatInput struct {
	UserID string `json:"userID"`
}

func CreateChat(w *response.Writer, r *http.Request, conn database.Connection) {
	var input CreateChatInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	// Check that user id exists
	_, err := conn.GetUser(input.UserID)
	if err != nil {
		errorResponse := fmt.Sprintf("userID %s is invalid", input.UserID)
		w.WriteString(http.StatusBadRequest, errorResponse)
		return
	}

	user, ok := getRequestContext(w, r, &struct{}{})
	if !ok {
		return
	}
	userIDs := [2]string{user.ID, input.UserID}

	// Check if chat already exists
	if chat := conn.GetChatByUserIDs(userIDs); chat != nil {
		w.WriteString(http.StatusConflict, "chat already exists")
		return
	}

	newChat := &database.Chat{
		ID:      uuid.NewString(),
		UserIDs: userIDs,
	}
	conn.SetChat(newChat)
	w.WriteJSON(http.StatusCreated, newChat)
}
