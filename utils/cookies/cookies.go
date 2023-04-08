package cookies

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/raphael-p/beango/utils/response"
)

type Cookie string

const (
	SESSION Cookie = "beango-session"
)

func Get(r *http.Request, name Cookie) (string, error) {
	cookies := r.Cookies()
	var matchingCookies []string
	for _, cookie := range cookies {
		if cookie.Name == string(name) {
			matchingCookies = append(matchingCookies, cookie.Value)
		}
	}
	numberOfMatches := len(matchingCookies)

	if numberOfMatches == 0 {
		return "", fmt.Errorf("no cookie found with the name %s", name)
	}
	if numberOfMatches > 1 {
		return "", fmt.Errorf("%d cookies found with the name %s", numberOfMatches, name)
	}
	return matchingCookies[0], nil
}

func Set(w *response.Writer, name Cookie, sessionID string, expiryDate time.Time) error {
	if name == "" {
		return errors.New("a cookie cannot have an empty name")
	}
	for _, cookieInResponse := range w.Header()["Set-Cookie"] {
		if strings.Contains(cookieInResponse, fmt.Sprint(string(name)+"=")) {
			return fmt.Errorf("response header already sets a cookie with the name %s", name)
		}
	}

	cookie := &http.Cookie{
		Name:     string(name),
		Value:    sessionID,
		Expires:  expiryDate,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	return nil
}

func Invalidate(w *response.Writer, name Cookie) error {
	return Set(w, name, "", time.Unix(0, 0))
}
