package auth

import (
	"errors"
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

type Connection interface {
	GetUser(id string) (*database.User, error)
	CheckSession(id string) (*database.Session, bool)
}

func Authentication(w *response.Writer, req *http.Request, conn Connection) (*http.Request, bool) {
	userID, err := getUserIDFromCookie(w, req, conn)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, false
	}
	user, err := conn.GetUser(userID)
	if err != nil {
		w.WriteString(http.StatusNotFound, "user not found during authentication")
		return nil, false
	}
	req, err = context.SetUser(req, user)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return nil, false
	}
	return req, true
}

func getUserIDFromCookie(w *response.Writer, req *http.Request, conn Connection) (string, error) {
	cookieName := cookies.SESSION
	sessionID, err := cookies.Get(req, cookieName)
	if err != nil {
		return "", err
	}
	session, ok := conn.CheckSession(sessionID)
	if !ok {
		err := cookies.Invalidate(w, cookieName)
		if err != nil {
			logger.Error(err.Error())
		}
		return "", errors.New("cookie or session is invalid")
	}
	return session.UserID, nil
}
