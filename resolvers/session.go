package resolvers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
	"golang.org/x/crypto/bcrypt"
)

type SessionInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateSession(w *response.Writer, r *http.Request) {
	sessionID, err := cookies.Get(r, cookies.SESSION)
	if err == nil {
		_, ok := database.CheckSession(sessionID)
		if ok {
			w.WriteString(http.StatusBadRequest, "there already is a valid session cookie in the request")
			return
		}
	}

	var input SessionInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	if user, err := database.GetUserByUsername(input.Username); err == nil {
		err := bcrypt.CompareHashAndPassword(user.Key, []byte(input.Password))
		if err == nil {
			sessionID := uuid.NewString()
			expiryDuration := time.Duration(config.Values.Session.SecondsUntilExpiry) * time.Second
			expiryDate := time.Now().UTC().Add(expiryDuration)
			err = cookies.Set(w, cookies.SESSION, sessionID, expiryDate)
			if err != nil {
				logger.Error(fmt.Sprint("session cookie creation failed: ", err))
				w.WriteString(http.StatusInternalServerError, "")
				return
			}
			database.SetSession(database.Session{
				ID:         sessionID,
				UserID:     user.ID,
				ExpiryDate: expiryDate,
			})
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
}
