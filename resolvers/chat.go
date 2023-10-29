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

func getChatsDatabase(userID int64, conn database.Connection) ([]getChatsOutput, *resolverutils.HTTPError) {
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

		if chat.Name == "" {
			chat.Name = generateChatName(userID, users)
		}
		chatOutput[i] = getChatsOutput{chat, stripUserFields(users...)}
	}
	return chatOutput, nil
}

func GetChats(w *response.Writer, r *http.Request, conn database.Connection) {
	user, _, httpError := resolverutils.GetRequestContext(r)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	chats, httpError := getChatsDatabase(user.ID, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	w.WriteJSON(http.StatusOK, chats)
}

type createPrivateChatInput struct {
	UserID int64 `json:"userID"`
}

func validateCreatePrivateChatInput(input *createPrivateChatInput, sessionUserID int64) *resolverutils.HTTPError {
	if sessionUserID == input.UserID {
		return &resolverutils.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "cannot create a chat with yourself",
		}
	}
	return nil
}

func createPrivateChatDatabase(sessionUserID, inputUserID int64, conn database.Connection) (*database.Chat, *resolverutils.HTTPError) {
	// Check that input user exists
	if user, _ := conn.GetUser(inputUserID); user == nil {
		return nil, &resolverutils.HTTPError{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("userID %d is invalid", inputUserID),
		}
	}

	userIDs := [2]int64{sessionUserID, inputUserID}

	// Check if chat already chatExists
	chatExists, err := conn.CheckPrivateChatExists(userIDs)
	if err != nil {
		return nil, resolverutils.HandleDatabaseError(err)
	}
	if chatExists {
		return nil, &resolverutils.HTTPError{
			Status:  http.StatusConflict,
			Message: "chat already exists",
		}
	}

	newChat := &database.Chat{Type: database.PRIVATE_CHAT}
	newChat, err = conn.SetChat(newChat, userIDs[:]...)
	if err != nil {
		return nil, resolverutils.HandleDatabaseError(err)
	}
	return newChat, nil
}

func CreatePrivateChat(w *response.Writer, r *http.Request, conn database.Connection) {
	var input createPrivateChatInput
	user, _, httpError := resolverutils.GetRequestBodyAndContext(r, &input)
	if resolverutils.ProcessHTTPError(w, httpError) ||
		resolverutils.ProcessHTTPError(w, validateCreatePrivateChatInput(&input, user.ID)) {
		return
	}

	newChat, httpError := createPrivateChatDatabase(user.ID, input.UserID, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}
	w.WriteJSON(http.StatusCreated, newChat)
}
