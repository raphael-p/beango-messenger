package resolvers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
)

func GetMessages(w *utils.ResponseWriter, r *http.Request) {
	_, vals := utils.MapValues(database.Messages)
	w.JSONResponse(http.StatusOK, vals)
}

type SendMessageInput struct {
	ChatId  string `json:"chatId"`
	Content string `json:"content"`
}

func SendMessage(w *utils.ResponseWriter, r *http.Request) {
	var input SendMessageInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	newMessage := database.Message{
		Id:      uuid.New().String(),
		ChatId:  input.ChatId,
		Content: input.Content,
	}
	database.Messages[newMessage.Id] = newMessage
	w.JSONResponse(http.StatusAccepted, newMessage)
}
