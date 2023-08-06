package resolvers

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/response"
	"github.com/raphael-p/beango/utils/validate"
	"golang.org/x/crypto/bcrypt"
)

type UserOutput struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

func stripFields(user *database.User) *UserOutput {
	return &UserOutput{user.ID, user.Username, user.DisplayName}
}

type CreateUserInput struct {
	Username    string                     `json:"username"`
	DisplayName validate.JSONField[string] `json:"displayName" optional:"true"`
	Password    string                     `json:"password"`
}

func CreateUser(w *response.Writer, r *http.Request, conn database.Connection) {
	var input CreateUserInput
	if ProcessHTTPError(w, getRequestBody(r, &input)) {
		return
	}

	if user, _ := conn.GetUserByUsername(input.Username); user != nil {
		w.WriteString(http.StatusConflict, "username is taken")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		w.WriteString(http.StatusBadRequest, err.Error())
		return
	}

	newUser := &database.User{
		Username:    input.Username,
		DisplayName: input.DisplayName.Value,
		Key:         hash,
	}
	if newUser.DisplayName == "" {
		newUser.DisplayName = input.Username
	}
	newUser, err = conn.SetUser(newUser)
	if err != nil {
		ProcessHTTPError(w, HandleDatabaseError(err))
		return
	}
	w.WriteJSON(http.StatusCreated, stripFields(newUser))
}

func GetUserByName(w *response.Writer, r *http.Request, conn database.Connection) {
	_, params, httpError := getRequestContext(r, USERNAME_KEY)
	if ProcessHTTPError(w, httpError) {
		return
	}

	user, _ := conn.GetUserByUsername(params[USERNAME_KEY])
	if user == nil {
		w.WriteString(http.StatusNotFound, "user not found")
		return
	}
	w.WriteJSON(http.StatusOK, stripFields(user))
}
