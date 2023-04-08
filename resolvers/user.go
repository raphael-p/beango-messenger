package resolvers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
	"golang.org/x/crypto/bcrypt"
)

type UserOutput struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

func stripFields(user *database.User) *UserOutput {
	return &UserOutput{user.ID, user.Username, user.DisplayName}
}

type CreateUserInput struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName" optional:"true"`
	Password    string `json:"password"`
}

func CreateUser(w *response.Writer, r *http.Request) {
	var input CreateUserInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	if user, _ := database.GetUserByUsername(input.Username); user != nil {
		w.WriteString(http.StatusConflict, "username is taken")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		w.WriteString(http.StatusBadRequest, err.Error())
		return
	}

	newUser := &database.User{
		ID:          uuid.New().String(),
		Username:    input.Username,
		DisplayName: input.DisplayName,
		Key:         hash,
	}
	if newUser.DisplayName == "" {
		newUser.DisplayName = input.Username
	}
	database.SetUser(newUser)
	w.WriteJSON(http.StatusCreated, stripFields(newUser))
}

func GetUserByName(w *response.Writer, r *http.Request) {
	paramKeys := []string{"username"}
	_, params, ok := getRequestContext(w, r, paramKeys...)
	if !ok {
		return
	}
	username := params[paramKeys[0]]

	user, _ := database.GetUserByUsername(username)
	if user == nil {
		w.WriteString(http.StatusNotFound, "user not found")
		return
	}
	w.WriteJSON(http.StatusOK, stripFields(user))
}
