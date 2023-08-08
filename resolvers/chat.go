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
func ChatsDatabase(userID int64, conn database.Connection) ([]GetChatsOutput, *HTTPError) {
	chats, err := conn.GetChatsByUserID(userID)
	if err != nil {
		return nil, HandleDatabaseError(err)
	}

	chatOutput := make([]GetChatsOutput, len(chats))
	for i, chat := range chats {
		users, err := conn.GetUsersByChatID(chat.ID)
		if err != nil {
			return nil, HandleDatabaseError(err)
		}

		outputUsers := make([]UserOutput, len(users))
		for j, user := range users {
			outputUsers[j] = *stripFields(&user)
		}

		if chat.Name == "" {
			chat.Name = generateChatName(userID, users)
		}
		chatOutput[i] = GetChatsOutput{chat, outputUsers}
	}
	return chatOutput, nil
}

func GetChats(w *response.Writer, r *http.Request, conn database.Connection) {
	user, _, httpError := getRequestContext(r)
	if ProcessHTTPError(w, httpError) {
		return
	}
	chats, httpError := ChatsDatabase(user.ID, conn)
	if ProcessHTTPError(w, httpError) {
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
	user, _, httpError := getRequestBodyAndContext(r, &input)
	if ProcessHTTPError(w, httpError) {
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
		ProcessHTTPError(w, HandleDatabaseError(err))
		return
	}
	if exists {
		w.WriteString(http.StatusConflict, "chat already exists")
		return
	}

	newChat := &database.Chat{Type: database.PRIVATE_CHAT}
	newChat, err = conn.SetChat(newChat, userIDs[:]...)
	if err != nil {
		ProcessHTTPError(w, HandleDatabaseError(err))
		return
	}
	w.WriteJSON(http.StatusCreated, newChat)
}
