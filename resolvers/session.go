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

var AUTH_COOKIE = "beango-session"

func CreateSession(w *utils.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(AUTH_COOKIE)
	if err == nil {
		session := database.GetSession(cookie.Value)
		if session != nil {
			w.StringResponse(http.StatusOK, "already authenticated")
			return
		}
	}

	var input SessionInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	user := database.GetUserByUsername(input.Username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.Key, []byte(input.Password))
		if err == nil {
			sessionId := uuid.New().String()
			expiryDate := time.Now().Add(24 * time.Hour)
			cookie := &http.Cookie{
				Name:     AUTH_COOKIE,
				Value:    sessionId,
				Expires:  expiryDate,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			}
			session := database.Session{
				Id:         sessionId,
				UserId:     user.Id,
				ExpiryDate: expiryDate,
			}
			database.AddSession(session)
			http.SetCookie(w, cookie)
			w.StringResponse(http.StatusOK, "authentication successful")
			return
		}
	}

	w.StringResponse(http.StatusUnauthorized, "authentication failed")
}
