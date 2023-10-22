package resolvers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/response"
)

type getChatsOutput struct {
	database.Chat
	Users []userOutput `json:"users"`
}

func generateChatName(userID int64, users []database.User) string {
	var displayNames []string
	for _, user := range users {
		if user.ID != userID {
			displayNames = append(displayNames, user.DisplayName)
		}
	}
	sort.Strings(displayNames)

	return strings.Join(displayNames, ", ")
}

func chatsDatabase(userID int64, conn database.Connection) ([]getChatsOutput, *resolverutils.HTTPError) {
	chats, err := conn.GetChatsByUserID(userID)
	if err != nil {
		return nil, resolverutils.HandleDatabaseError(err)
	}

	chatOutput := make([]getChatsOutput, len(chats))
	for i, chat := range chats {
		users, err := conn.GetUsersByChatID(chat.ID)
		if err != nil {
			return nil, resolverutils.HandleDatabaseError(err)
		}

		outputUsers := make([]userOutput, len(users))
		for j, user := range users {
			outputUsers[j] = *stripFields(&user)
		}

		if chat.Name == "" {
			chat.Name = generateChatName(userID, users)
		}
		chatOutput[i] = getChatsOutput{chat, outputUsers}
	}
	return chatOutput, nil
}

func GetChats(w *response.Writer, r *http.Request, conn database.Connection) {
	user, _, httpError := resolverutils.GetRequestContext(r)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	chats, httpError := chatsDatabase(user.ID, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	w.WriteJSON(http.StatusOK, chats)
}

type createChatInput struct {
	UserID int64 `json:"userID"`
}

func CreatePrivateChat(w *response.Writer, r *http.Request, conn database.Connection) {
	var input createChatInput
	user, _, httpError := resolverutils.GetRequestBodyAndContext(r, &input)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	if user.ID == input.UserID {
		w.WriteString(http.StatusBadRequest, "cannot create a chat with yourself")
		return
	}
	userIDs := [2]int64{user.ID, input.UserID}

	// Check that user id exists
	if user, _ := conn.GetUser(input.UserID); user == nil {
		errorResponse := fmt.Sprintf("userID %d is invalid", input.UserID)
		w.WriteString(http.StatusBadRequest, errorResponse)
		return
	}

	// Check if chat already chatExists
	chatExists, err := conn.CheckPrivateChatExists(userIDs)
	if err != nil {
		resolverutils.ProcessHTTPError(w, resolverutils.HandleDatabaseError(err))
		return
	}
	if chatExists {
		w.WriteString(http.StatusConflict, "chat already exists")
		return
	}

	newChat := &database.Chat{Type: database.PRIVATE_CHAT}
	newChat, err = conn.SetChat(newChat, userIDs[:]...)
	if err != nil {
		resolverutils.ProcessHTTPError(w, resolverutils.HandleDatabaseError(err))
		return
	}
	w.WriteJSON(http.StatusCreated, newChat)
}
