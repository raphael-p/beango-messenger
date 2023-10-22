package resolvers

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/response"
	"github.com/raphael-p/beango/utils/validate"
	"golang.org/x/crypto/bcrypt"
)

type userOutput struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

func stripFields(user *database.User) *userOutput {
	return &userOutput{user.ID, user.Username, user.DisplayName}
}

type createUserInput struct {
	Username    string                     `json:"username"`
	DisplayName validate.JSONField[string] `json:"displayName" optional:"true"`
	Password    string                     `json:"password"`
}

func createUserDatabase(username, displayName, password string, conn database.Connection) (*userOutput, *resolverutils.HTTPError) {
	if user, _ := conn.GetUserByUsername(username); user != nil {
		return nil, &resolverutils.HTTPError{Status: http.StatusConflict, Message: "username is taken"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, &resolverutils.HTTPError{Status: http.StatusBadRequest, Message: err.Error()}
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
		return nil, resolverutils.HandleDatabaseError(err)
	}

	return stripFields(newUser), nil
}

func CreateUser(w *response.Writer, r *http.Request, conn database.Connection) {
	var input createUserInput
	if resolverutils.ProcessHTTPError(w, resolverutils.GetRequestBody(r, &input)) {
		return
	}

	newUser, httpError := createUserDatabase(input.Username, input.DisplayName.Value, input.Password, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	w.WriteJSON(http.StatusCreated, newUser)
}

func GetUserByName(w *response.Writer, r *http.Request, conn database.Connection) {
	_, params, httpError := resolverutils.GetRequestContext(r, resolverutils.USERNAME_KEY)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	user, _ := conn.GetUserByUsername(params.Username)
	if user == nil {
		w.WriteString(http.StatusNotFound, "user not found")
		return
	}
	w.WriteJSON(http.StatusOK, stripFields(user))
}
