package resolvers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
)

func GetChats(c *gin.Context) {
	_, vals := utils.MapValues(database.Chats)
	c.IndentedJSON(http.StatusOK, vals)
}

func CreateChat(c *gin.Context) {
	var newChat database.Chat

	if err := c.BindJSON(&newChat); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "malformed request body")
		return
	}

	if len(newChat.UserIds) != 2 {
		c.IndentedJSON(http.StatusBadRequest, "chat must have exactly two users")
		return
	}

	// Check that user ids exist
	var missingIds []string
	for _, userId := range newChat.UserIds {
		_, exists := database.Users[userId]
		if !exists {
			missingIds = append(missingIds, userId)
		}
	}
	if len(missingIds) != 0 {
		c.IndentedJSON(http.StatusBadRequest, "invalid userId(s): "+strings.Join(missingIds, ", "))
		return
	}

	// Check if chat already exists
	var chatId *string
	for key, value := range database.Chats {
		if (newChat.UserIds[0] == value.UserIds[0] &&
			newChat.UserIds[1] == value.UserIds[1]) ||
			(newChat.UserIds[0] == value.UserIds[1] &&
				newChat.UserIds[1] == value.UserIds[0]) {
			chatId = &key
		}
	}
	if chatId != nil {
		newChat.Id = *chatId
	} else {
		newChat.Id = uuid.New().String()
	}
	database.Chats[newChat.Id] = newChat
	c.IndentedJSON(http.StatusCreated, newChat)
}
