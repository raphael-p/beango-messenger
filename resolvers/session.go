package resolvers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/httputils"
	"golang.org/x/crypto/bcrypt"
)

type SessionInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateSession(w *httputils.ResponseWriter, r *http.Request) {
	sessionId, _ := httputils.GetCookieValue(httputils.AUTH_COOKIE, r)
	if sessionId != "" {
		_, ok := database.CheckSession(sessionId)
		if ok {
			w.StringResponse(http.StatusBadRequest, "there already is a valid session cookie in the request")
			return
		}
	}

	var input SessionInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	user, _ := database.GetUserByUsername(input.Username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.Key, []byte(input.Password))
		if err == nil {
			sessionId := uuid.NewString()
			expiryDuration := time.Duration(config.Values.Session.SecondsUntilExpiry) * time.Second
			expiryDate := time.Now().Add(expiryDuration)
			httputils.SetCookie(httputils.AUTH_COOKIE, sessionId, expiryDate, w)
			database.SetSession(database.Session{
				Id:         sessionId,
				UserId:     user.Id,
				ExpiryDate: expiryDate,
			})
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
}
