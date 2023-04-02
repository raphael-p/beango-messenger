package resolvers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/response"
	"golang.org/x/crypto/bcrypt"
)

type SessionInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateSession(w *response.Writer, r *http.Request) {
	sessionId, _ := cookies.Get(r, cookies.SESSION)
	if sessionId != "" {
		_, ok := database.CheckSession(sessionId)
		if ok {
			w.WriteString(http.StatusBadRequest, "there already is a valid session cookie in the request")
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
			cookies.Set(w, cookies.SESSION, sessionId, expiryDate)
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
