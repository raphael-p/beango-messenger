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

func createUserDatabase(username, displayName, password string, conn database.Connection) (*UserOutput, *HTTPError) {
	if user, _ := conn.GetUserByUsername(username); user != nil {
		return nil, &HTTPError{http.StatusConflict, "username is taken"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, &HTTPError{http.StatusBadRequest, err.Error()}
	}

	newUser := &database.User{
		Username:    username,
		DisplayName: displayName,
		Key:         hash,
	}
	if newUser.DisplayName == "" {
		newUser.DisplayName = username
	}
	newUser, err = conn.SetUser(newUser)
	if err != nil {
		return nil, HandleDatabaseError(err)
	}

	return stripFields(newUser), nil
}

func CreateUser(w *response.Writer, r *http.Request, conn database.Connection) {
	var input CreateUserInput
	if ProcessHTTPError(w, getRequestBody(r, &input)) {
		return
	}

	newUser, httpError := createUserDatabase(input.Username, input.DisplayName.Value, input.Password, conn)
	if ProcessHTTPError(w, httpError) {
		return
	}

	w.WriteJSON(http.StatusCreated, newUser)
}

func GetUserByName(w *response.Writer, r *http.Request, conn database.Connection) {
	_, params, httpError := getRequestContext(r, USERNAME_KEY)
	if ProcessHTTPError(w, httpError) {
		return
	}

	user, _ := conn.GetUserByUsername(params.Username)
	if user == nil {
		w.WriteString(http.StatusNotFound, "user not found")
		return
	}
	w.WriteJSON(http.StatusOK, stripFields(user))
}
