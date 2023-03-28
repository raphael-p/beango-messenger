package resolvers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/httputils"
	"golang.org/x/crypto/bcrypt"
)

type UserOutput struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

func stripFields(user *database.User) *UserOutput {
	return &UserOutput{user.Id, user.Username, user.DisplayName}
}

type CreateUserInput struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName" optional:"true"`
	Password    string `json:"password"`
}

func CreateUser(w *httputils.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	if user, _ := database.GetUserByUsername(input.Username); user != nil {
		w.StringResponse(http.StatusConflict, "username is taken")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		w.StringResponse(http.StatusBadRequest, err.Error())
		return
	}

	newUser := &database.User{
		Id:          uuid.New().String(),
		Username:    input.Username,
		DisplayName: input.DisplayName,
		Key:         hash,
	}
	if newUser.DisplayName == "" {
		newUser.DisplayName = input.Username
	}
	database.SetUser(newUser)
	w.JSONResponse(http.StatusCreated, stripFields(newUser))
}

func GetUserByName(w *httputils.ResponseWriter, r *http.Request) {
	username, err := httputils.GetParamFromContext(r, "username")
	if err != nil {
		w.StringResponse(http.StatusInternalServerError, err.Error())
	}
	user, _ := database.GetUserByUsername(username)
	if user == nil {
		w.StringResponse(http.StatusNotFound, "user not found")
		return
	}
	w.JSONResponse(http.StatusOK, stripFields(user))
}
