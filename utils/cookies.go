package utils

import (
	"net/http"
	"time"
)

type Cookie string

const (
	AUTH_COOKIE Cookie = "beango-session"
)

func GetCookieValue(name Cookie, r *http.Request) (string, error) {
	cookie, err := r.Cookie(string(name))
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func SetCookie(name Cookie, sessionId string, expiryDate time.Time, w *ResponseWriter) {
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

func InvalidateCookie(cookie Cookie, w *ResponseWriter) {
	invalidCookie := &http.Cookie{
		Name:    string(cookie),
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	}
	w.Header().Set("Set-Cookie", invalidCookie.String()+"; SameSite=Strict; Secure")
}
