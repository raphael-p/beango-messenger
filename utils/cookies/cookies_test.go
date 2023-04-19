package cookies

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/response"
)

func findCookies(w *response.Writer, names ...string) []string {
	var matches []string
	for _, cookieInResponse := range w.Header()["Set-Cookie"] {
		for _, name := range names {
			if strings.Contains(cookieInResponse, fmt.Sprint(name+"=")) {
				matches = append(matches, cookieInResponse)
			}
		}
	}
	return matches
}

func makeCookie(name, value string, expiry time.Time) string {
	return fmt.Sprintf(
		"%s=%s; Path=/; Expires=%s; HttpOnly; Secure; SameSite=Strict",
		name,
		value,
		expiry.Format("Mon, 02 Jan 2006 15:04:05 GMT"),
	)
}

func TestGet_Single(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	name := "my-session"
	xSessionID := "session-id"
	cookie := &http.Cookie{Name: name, Value: xSessionID}
	req.AddCookie(cookie)

	sessionID, err := Get(req, Cookie(name))
	assert.IsNil(t, err)
	assert.Equals(t, sessionID, xSessionID)
}

func TestGet_DifferentNames(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	name1 := "name-1"
	name2 := "name-2"
	xValue1 := "value-1"
	xValue2 := "value-2"
	cookie1 := &http.Cookie{Name: name1, Value: xValue1}
	cookie2 := &http.Cookie{Name: name2, Value: xValue2}
	req.AddCookie(cookie1)
	req.AddCookie(cookie2)

	value1, err1 := Get(req, Cookie(name1))
	value2, err2 := Get(req, Cookie(name2))
	assert.IsNil(t, err1)
	assert.IsNil(t, err2)
	assert.Equals(t, value1, xValue1)
	assert.Equals(t, value2, xValue2)
}

func TestGet_SameName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	name := "name"
	cookie1 := &http.Cookie{Name: name, Value: "value-1"}
	cookie2 := &http.Cookie{Name: name, Value: "value-2"}
	req.AddCookie(cookie1)
	req.AddCookie(cookie2)

	_, err1 := Get(req, Cookie(name))
	_, err2 := Get(req, Cookie(name))
	xErrorMessage := fmt.Sprint("2 cookies found with the name ", name)
	assert.ErrorHasMessage(t, err1, xErrorMessage)
	assert.ErrorHasMessage(t, err2, xErrorMessage)
}

func TestGet_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{
		Name:    "test-cookie",
		Value:   "test-session-id",
		Expires: time.Now().UTC().Add(5 * time.Second),
		Path:    "/",
	}
	req.AddCookie(cookie)

	_, err := Get(req, SESSION)
	assert.ErrorHasMessage(t, err, fmt.Sprint("no cookie found with the name ", SESSION))
}

func TestSet_Single(t *testing.T) {
	w := mocks.MakeResponseWriter()
	name := "test-name"
	value := "test-value"
	expiry := time.Now().UTC().Add(24 * time.Hour)

	err := Set(w, Cookie(name), value, expiry)
	assert.IsNil(t, err)
	cookies := findCookies(w, name)
	xCookies := []string{makeCookie(name, value, expiry)}
	assert.DeepEquals(t, cookies, xCookies)
}

func TestSet_EmptyName(t *testing.T) {
	w := mocks.MakeResponseWriter()
	name := ""
	value := "test-value"
	expiry := time.Now().UTC().Add(24 * time.Hour)

	err := Set(w, Cookie(name), value, expiry)
	assert.ErrorHasMessage(t, err, "a cookie cannot have an empty name")
}

func TestSet_DifferentNames(t *testing.T) {
	w := mocks.MakeResponseWriter()
	name1 := "test-name-1"
	value1 := "test-value-1"
	name2 := "test-name-2"
	value2 := "test-value-2"
	expiry := time.Now().UTC().Add(24 * time.Hour)

	err := Set(w, Cookie(name1), value1, expiry)
	assert.IsNil(t, err)
	err = Set(w, Cookie(name2), value2, expiry)
	assert.IsNil(t, err)
	cookies := findCookies(w, name1, name2)
	xCookies := []string{
		makeCookie(name1, value1, expiry),
		makeCookie(name2, value2, expiry),
	}
	assert.DeepEquals(t, cookies, xCookies)
}

func TestSet_SameName(t *testing.T) {
	w := mocks.MakeResponseWriter()
	name := "test-name"
	value := "test-value"
	expiry := time.Now().UTC().Add(24 * time.Hour)

	err := Set(w, Cookie(name), value, expiry)
	assert.IsNil(t, err)
	err = Set(w, Cookie(name), value, expiry)
	xErrorMessage := fmt.Sprint("response header already sets a cookie with the name ", name)
	assert.ErrorHasMessage(t, err, xErrorMessage)
}

func TestInvalidate(t *testing.T) {
	w := mocks.MakeResponseWriter()
	name := "test-name"
	Invalidate(w, Cookie(name))

	cookie := w.Header().Get("Set-Cookie")
	xCookie := makeCookie(name, "", time.Unix(0, 0).UTC())
	assert.Equals(t, cookie, xCookie)
}
