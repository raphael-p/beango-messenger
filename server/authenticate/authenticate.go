package authenticate

import (
	"errors"
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func FromCookie(w *response.Writer, req *http.Request, conn database.Connection) (*http.Request, bool) {
	userID, err := getUserIDFromCookie(w, req, conn)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return req, false
	}
	user, err := conn.GetUser(userID)
	if err != nil {
		w.WriteString(http.StatusNotFound, "user not found during authentication")
		return req, false
	}
	req, err = context.SetUser(req, user)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return req, false
	}
	return req, true
}

func getUserIDFromCookie(w *response.Writer, req *http.Request, conn database.Connection) (int64, error) {
	cookieName := cookies.SESSION
	sessionID, err := cookies.Get(req, cookieName)
	if err != nil {
		return 0, err
	}
	session, ok := conn.CheckSession(sessionID)
	if !ok {
		err := cookies.Invalidate(w, cookieName)
		if err != nil {
			logger.Error(err.Error())
		}
		return 0, errors.New("cookie or session is invalid")
	}
	return session.UserID, nil
}
