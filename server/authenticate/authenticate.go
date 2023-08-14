package authenticate

import (
	"errors"
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func Auth(w *response.Writer, r *http.Request, conn database.Connection) (*http.Request, *resolvers.HTTPError) {
	userID, err := getUserIDFromCookie(w, r, conn)
	if err != nil {
		return r, &resolvers.HTTPError{Status: http.StatusUnauthorized}
	}

	user, err := conn.GetUser(userID)
	if user == nil {
		var httpError *resolvers.HTTPError
		if err != nil {
			httpError = resolvers.HandleDatabaseError(err)
		} else {
			httpError = &resolvers.HTTPError{
				Status:  http.StatusNotFound,
				Message: "user not found during authentication",
			}
		}
		return r, httpError
	}

	r, err = context.SetUser(r, user)
	if err != nil {
		logger.Error(err.Error())
		return r, &resolvers.HTTPError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return r, nil
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
