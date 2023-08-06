package resolvers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
)

// TODO: Use ChatExtended struct?
type GetChatsOutput struct {
	database.Chat
	Users []UserOutput `json:"users"`
}

// TODO: test
func generateChatName(userID int64, users []database.User) string {
	var displayNames []string
	for _, user := range users {
		if user.ID != userID {
			displayNames = append(displayNames, user.DisplayName)
		}
	}

	return strings.Join(displayNames, ", ")
}

// TODO: test
func GetChatsData(w *response.Writer, r *http.Request, conn database.Connection) ([]GetChatsOutput, bool) {
	user, _, ok := getRequestContext(w, r)
	if !ok {
		return nil, false
	}

	chats, err := conn.GetChatsByUserID(user.ID)
	if err != nil {
		ProcessHTTPError(HandleDatabaseError(err), w)
		return nil, false
	}

	chatOutput := make([]GetChatsOutput, len(chats))
	for i, chat := range chats {
		users, err := conn.GetUsersByChatID(chat.ID)
		if err != nil {
			ProcessHTTPError(HandleDatabaseError(err), w)
			return nil, false
		}

		outputUsers := make([]UserOutput, len(users))
		for j, user := range users {
			outputUsers[j] = *stripFields(&user)
		}

		if chat.Name == "" {
			chat.Name = generateChatName(user.ID, users)
		}
		chatOutput[i] = GetChatsOutput{chat, outputUsers}
	}
	return chatOutput, true
}

func GetChats(w *response.Writer, r *http.Request, conn database.Connection) {
	if chats, ok := GetChatsData(w, r, conn); ok {
		w.WriteJSON(http.StatusOK, chats)
	}
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
		ProcessHTTPError(HandleDatabaseError(err), w)
		return
	}
	if exists {
		w.WriteString(http.StatusConflict, "chat already exists")
		return
	}

	newChat := &database.Chat{Type: database.PRIVATE_CHAT}
	newChat, err = conn.SetChat(newChat, userIDs[:]...)
	if err != nil {
		ProcessHTTPError(HandleDatabaseError(err), w)
		return
	}
	w.WriteJSON(http.StatusCreated, newChat)
}
