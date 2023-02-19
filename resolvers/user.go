package resolvers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserOutput struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func stripFields(user database.User) UserOutput {
	return UserOutput{user.Id, user.Username}
}

type GetUsersOutput []UserOutput

func GetUsers(w *utils.ResponseWriter, r *http.Request) {
	_, vals := utils.MapValues(database.Users)

	var output GetUsersOutput
	for _, val := range vals {
		output = append(output, stripFields(val))
	}

	w.JSONResponse(http.StatusOK, output)
}

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateUser(w *utils.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if ok := bindRequestJSON(w, r, &input); !ok {
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

	newUser := database.User{
		Id:       uuid.New().String(),
		Username: input.Username,
		Key:      string(hash),
	}
	database.Users[newUser.Id] = newUser
	w.JSONResponse(http.StatusCreated, stripFields(newUser))
}
