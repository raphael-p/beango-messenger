package utils

import (
	"net/http"
	"time"
)

type Cookie string

const (
	AUTH_COOKIE Cookie = "beango-session"
)

func GetCookie(name Cookie, r *http.Request) *http.Cookie {
	cookie, err := r.Cookie(string(name))
	if err != nil {
		return nil
	}
	return cookie
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
