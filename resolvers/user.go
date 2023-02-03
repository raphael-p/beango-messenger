package resolvers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUsers(c *gin.Context) {
	_, vals := utils.MapValues(database.Users)
	c.IndentedJSON(http.StatusOK, vals)
}

func CreateUser(c *gin.Context) {
	var input CreateUserInput
	newUser := database.User{}

	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "malformed request body")
		return
	}

	if input.Username == "" {
		c.IndentedJSON(http.StatusBadRequest, "username is missing")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	newUser.Id = uuid.New().String()
	newUser.Username = input.Username
	newUser.Key = string(hash)
	database.Users[newUser.Id] = newUser
	c.IndentedJSON(http.StatusCreated, newUser)
}
