package resolvers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
	"golang.org/x/crypto/bcrypt"
)

type SessionInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateSession(w *utils.ResponseWriter, r *http.Request) {
	sessionId, _ := utils.GetCookieValue(utils.AUTH_COOKIE, r)
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
			expiryDate := time.Now().Add(24 * time.Hour)
			utils.SetCookie(utils.AUTH_COOKIE, sessionId, expiryDate, w)
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
