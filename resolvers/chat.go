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
	chats, err := conn.GetChatsByUserID(user.ID)
	if err != nil {
		HandleDatabaseError(w, err)
		return
	}
	w.WriteJSON(http.StatusOK, chats)
}

type CreateChatInput struct {
	UserID int64 `json:"userID"`
}

func CreatePrivateChat(w *response.Writer, r *http.Request, conn database.Connection) {
	// TODO: throw error if attempting to create self chat, create new chat type for that, called "note"
	var input CreateChatInput
	user, _, ok := getRequestBodyAndContext(w, r, &input)
	if !ok {
		return
	}
	userIDs := [2]int64{user.ID, input.UserID}

	// Check that user id exists
	if user, _ := conn.GetUser(input.UserID); user == nil {
		errorResponse := fmt.Sprintf("userID %d is invalid", input.UserID)
		w.WriteString(http.StatusBadRequest, errorResponse)
		return
	}

	// Check if chat already exists
	exists, err := conn.CheckPrivateChatExists(userIDs)
	if err != nil {
		HandleDatabaseError(w, err)
		return
	}
	if exists {
		w.WriteString(http.StatusConflict, "chat already exists")
		return
	}

	newChat := &database.Chat{Type: database.PRIVATE_CHAT}
	newChat, err = conn.SetChat(newChat, userIDs[:]...)
	if err != nil {
		HandleDatabaseError(w, err)
		return
	}
	w.WriteJSON(http.StatusCreated, newChat)
}
