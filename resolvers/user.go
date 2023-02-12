package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUsers(w *utils.ResponseWriter, r *http.Request) {
	_, vals := utils.MapValues(database.Users)
	w.JSONResponse(http.StatusOK, vals)
}

func CreateUser(w *utils.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	newUser := database.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		fmt.Println(err.Error())
		w.StringResponse(http.StatusBadRequest, "malformed request body")
		return
	}

	if input.Username == "" {
		w.StringResponse(http.StatusBadRequest, "username is missing")
		return
	}

	for _, value := range database.Users {
		if value.Username == input.Username {
			w.StringResponse(http.StatusConflict, "username is taken")
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		w.StringResponse(http.StatusBadRequest, err.Error())
		return
	}
	newUser.Id = uuid.New().String()
	newUser.Username = input.Username
	newUser.Key = string(hash)
	database.Users[newUser.Id] = newUser
	w.JSONResponse(http.StatusCreated, newUser)
}
