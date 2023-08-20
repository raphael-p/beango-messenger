package resolvers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
	"golang.org/x/crypto/bcrypt"
)

type SessionInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TODO unit test
func checkCredentials(username, password string, conn database.Connection) (int64, *resolverutils.HTTPError) {
	unauthorised := func() *resolverutils.HTTPError {
		return &resolverutils.HTTPError{http.StatusUnauthorized, "login credentials are incorrect"}
	}
	user, _ := conn.GetUserByUsername(username)
	if user == nil {
		return 0, unauthorised()
	}
	err := bcrypt.CompareHashAndPassword(user.Key, []byte(password))
	if err != nil {
		return 0, unauthorised()
	}

	return user.ID, nil
}

// TODO unit test
func setSession(w *response.Writer, userID int64, conn database.Connection) *resolverutils.HTTPError {
	expiryDuration := time.Duration(config.Values.Session.SecondsUntilExpiry) * time.Second
	expiryDate := time.Now().UTC().Add(expiryDuration)
	sessionID := uuid.NewString()
	if err := cookies.Set(w, cookies.SESSION, sessionID, expiryDate); err != nil {
		logger.Error(fmt.Sprint("failed to create session cookie: ", err))
		return &resolverutils.HTTPError{http.StatusInternalServerError, "failed to create session cookie"}
	}
	conn.SetSession(database.Session{
		ID:         sessionID,
		UserID:     userID,
		ExpiryDate: expiryDate,
	})
	return nil
}

func CreateSession(w *response.Writer, r *http.Request, conn database.Connection) {
	if sessionID, err := cookies.Get(r, cookies.SESSION); err == nil {
		if _, ok := conn.CheckSession(sessionID); ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	var input SessionInput
	if resolverutils.ProcessHTTPError(w, resolverutils.GetRequestBody(r, &input)) {
		return
	}

	userID, httpError := checkCredentials(input.Username, input.Password, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	if resolverutils.ProcessHTTPError(w, setSession(w, userID, conn)) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
