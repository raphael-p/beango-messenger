package resolvers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
)

func GetMessages(c *gin.Context) {
	_, vals := utils.MapValues(database.Messages)
	c.IndentedJSON(http.StatusOK, vals)
}

func SendMessage(c *gin.Context) {
	var newMessage database.Message

	if err := c.BindJSON(&newMessage); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "malformed request body")
		return
	}

	newMessage.Id = uuid.New().String()
	database.Messages[newMessage.Id] = newMessage
	c.IndentedJSON(http.StatusAccepted, newMessage)
}
