package resolvers

import (
	"fmt"
	"net/http"

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
	UserID int `json:"userID"`
}

func CreatePrivateChat(w *response.Writer, r *http.Request, conn database.Connection) {
	var input CreateChatInput
	user, _, ok := getRequestBodyAndContext(w, r, &input)
	if !ok {
		return
	}
	userIDs := [2]int{user.ID, input.UserID}

	// Check that user id exists
	_, err := conn.GetUser(input.UserID)
	if err != nil {
		errorResponse := fmt.Sprintf("userID %d is invalid", input.UserID)
		w.WriteString(http.StatusBadRequest, errorResponse)
		return
	}

	// Check if chat already exists
	if conn.CheckPrivateChatExists(userIDs) {
		w.WriteString(http.StatusConflict, "chat already exists")
		return
	}

	newChat := &database.Chat{ChatType: database.PRIVATE_CHAT}
	conn.SetChat(newChat, userIDs[:]...)
	w.WriteJSON(http.StatusCreated, newChat)
}
