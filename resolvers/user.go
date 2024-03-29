package resolvers

import (
	"net/http"
	"regexp"
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

func stripUserFields(users ...database.User) []userOutput {
	output := make([]userOutput, len(users))
	for idx, user := range users {
		output[idx] = userOutput{user.ID, user.Username, user.DisplayName}
	}
	return output
}

type createUserInput struct {
	Username    string                     `json:"username"`
	DisplayName validate.JSONField[string] `json:"displayName" optional:"true"`
	Password    string                     `json:"password"`
}

func validateCreateUserInput(input *createUserInput) *resolverutils.HTTPError {
	input.Username = strings.TrimSpace(input.Username)
	input.DisplayName.Value = strings.TrimSpace(input.DisplayName.Value)

	if len([]rune(input.Username)) > 15 {
		return &resolverutils.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "username must be shorter than 16 characters",
		}
	}
	if !regexp.MustCompile("^[a-zA-Z0-9_.]*$").MatchString(input.Username) {
		return &resolverutils.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "username may only contain alphanumeric characters and '_.'",
		}
	}

	if len([]rune(input.DisplayName.Value)) > 15 {
		return &resolverutils.HTTPError{
			Status:  http.StatusBadRequest,
			Message: "display name must be shorter than 16 characters",
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

	return &stripUserFields(*newUser)[0], nil
}

func CreateUser(w *response.Writer, r *http.Request, conn database.Connection) {
	var input createUserInput
	if resolverutils.ProcessHTTPError(w, resolverutils.GetRequestBody(r, &input)) {
		return
	}
	if resolverutils.ProcessHTTPError(w, validateCreateUserInput(&input)) {
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
	w.WriteJSON(http.StatusOK, stripUserFields(*user)[0])
}
