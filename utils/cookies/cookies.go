package cookies

import (
	"net/http"
	"time"

	"github.com/raphael-p/beango/utils/response"
)

type Cookie string

const (
	SESSION Cookie = "beango-session"
)

func Get(r *http.Request, name Cookie) (string, error) {
	cookie, err := r.Cookie(string(name))
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func Set(w *response.Writer, name Cookie, sessionId string, expiryDate time.Time) {
	cookie := &http.Cookie{
		Name:     string(name),
		Value:    sessionId,
		Expires:  expiryDate,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func Invalidate(w *response.Writer, cookie Cookie) {
	invalidCookie := &http.Cookie{
		Name:    string(cookie),
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	}
	w.Header().Set("Set-Cookie", invalidCookie.String()+"; SameSite=Strict; Secure")
}
