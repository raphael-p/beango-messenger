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

func CreateSession(w *response.Writer, r *http.Request, conn database.Connection) {
	if sessionID, err := cookies.Get(r, cookies.SESSION); err == nil {
		_, ok := conn.CheckSession(sessionID)
		if ok {
			w.WriteString(http.StatusBadRequest, "there already is a valid session cookie in the request")
			return
		}
	}

	var input SessionInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	user, err := conn.GetUserByUsername(input.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Key, []byte(input.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID := uuid.NewString()
	expiryDuration := time.Duration(config.Values.Session.SecondsUntilExpiry) * time.Second
	expiryDate := time.Now().UTC().Add(expiryDuration)
	err = cookies.Set(w, cookies.SESSION, sessionID, expiryDate)
	if err != nil {
		logger.Error(fmt.Sprint("failed to create session cookie: ", err))
		w.WriteString(http.StatusInternalServerError, "failed to create session cookie")
		return
	}
	conn.SetSession(database.Session{
		ID:         sessionID,
		UserID:     user.ID,
		ExpiryDate: expiryDate,
	})
	w.WriteHeader(http.StatusNoContent)
}
