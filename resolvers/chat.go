package resolvers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
)

func GetChats(w *utils.ResponseWriter, r *http.Request) {
	_, vals := utils.MapValues(database.Chats)
	w.JSONResponse(http.StatusOK, vals)
}

type CreateChatInput struct {
	UserIds []string `json:"userIds"`
}

func CreateChat(w *utils.ResponseWriter, r *http.Request) {
	var input CreateChatInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	if len(input.UserIds) != 2 {
		w.StringResponse(http.StatusBadRequest, "chat must have exactly two users")
		return
	}

	// Check that user ids exist
	var missingIds []string
	for _, userId := range input.UserIds {
		_, exists := database.Users[userId]
		if !exists {
			missingIds = append(missingIds, userId)
		}
	}
	if len(missingIds) != 0 {
		err := "invalid userId(s): " + strings.Join(missingIds, ", ")
		w.StringResponse(http.StatusBadRequest, err)
		return
	}
	newChat := database.Chat{UserIds: input.UserIds}

	// Check if chat already exists
	for key, value := range database.Chats {
		if (input.UserIds[0] == value.UserIds[0] &&
			input.UserIds[1] == value.UserIds[1]) ||
			(input.UserIds[0] == value.UserIds[1] &&
				input.UserIds[1] == value.UserIds[0]) {
			newChat.Id = key
		}
	}
	if newChat.Id == "" {
		newChat.Id = uuid.New().String()
		database.Chats[newChat.Id] = newChat
	}
	w.JSONResponse(http.StatusCreated, newChat)
}
