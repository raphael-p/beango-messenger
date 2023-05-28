package resolvers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
)

func GetChats(w *response.Writer, r *http.Request, conn database.Connection) {
	user, _, ok := getRequestContext(w, r)
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
	user, _, ok := getRequestBodyAndContext(w, r, &input)
	if !ok {
		return
	}
	userIDs := [2]string{user.ID, input.UserID}

	// Check that user id exists
	_, err := conn.GetUser(input.UserID)
	if err != nil {
		errorResponse := fmt.Sprintf("userID %s is invalid", input.UserID)
		w.WriteString(http.StatusBadRequest, errorResponse)
		return
	}

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
