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

type sessionInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func checkCredentials(username, password string, conn database.Connection) (int64, *resolverutils.HTTPError) {
	unauthorised := func() *resolverutils.HTTPError {
		return &resolverutils.HTTPError{
			Status:  http.StatusUnauthorized,
			Message: "login credentials are incorrect",
		}
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

func makeSession(userID int64) *database.Session {
	sessionID := uuid.NewString()
	expiryDuration := time.Duration(config.Values.Session.SecondsUntilExpiry) * time.Second
	expiryDate := time.Now().UTC().Add(expiryDuration)
	return &database.Session{
		ID:         sessionID,
		UserID:     userID,
		ExpiryDate: expiryDate,
	}
}

func setSession(w *response.Writer, session *database.Session, conn database.Connection) *resolverutils.HTTPError {
	if err := cookies.Set(w, cookies.SESSION, session.ID, session.ExpiryDate); err != nil {
		logger.Error(fmt.Sprint("failed to create session cookie: ", err))
		return &resolverutils.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "failed to create session cookie",
		}
	}
	conn.SetSession(*session)
	return nil
}

func CreateSession(w *response.Writer, r *http.Request, conn database.Connection) {
	if sessionID, err := cookies.Get(r, cookies.SESSION); err == nil {
		if _, ok := conn.CheckSession(sessionID); ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	var input sessionInput
	if resolverutils.ProcessHTTPError(w, resolverutils.GetRequestBody(r, &input)) {
		return
	}

	userID, httpError := checkCredentials(input.Username, input.Password, conn)
	if resolverutils.ProcessHTTPError(w, httpError) {
		return
	}

	if resolverutils.ProcessHTTPError(w, setSession(w, makeSession(userID), conn)) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
