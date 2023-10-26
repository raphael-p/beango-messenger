package resolvers

import (
	"net/http"
	"strings"

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

func createUserInputValidation(input *createUserInput) *resolverutils.HTTPError {
	input.Username = strings.TrimSpace(input.Username)
	input.DisplayName.Value = strings.TrimSpace(input.DisplayName.Value)

	if len([]rune(input.Username)) > 25 {
		return &resolverutils.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "username must be shorter than 26 characters",
		}
	}
	if strings.ContainsAny(input.Username, " \t\n\r") {
		return &resolverutils.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "username may not contain any spaces, tabs, or new lines",
		}
	}

	if len([]rune(input.DisplayName.Value)) > 25 {
		return &resolverutils.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "display name must be shorter than 26 characters",
		}
	}
	if strings.ContainsAny(input.DisplayName.Value, "\t\n\r") {
		return &resolverutils.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "display name may not contain any tabs or newlines",
		}
	}
	return nil
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
	if resolverutils.ProcessHTTPError(w, createUserInputValidation(&input)) {
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
