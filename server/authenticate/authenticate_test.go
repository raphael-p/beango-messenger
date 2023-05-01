package authenticate

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func setup(name, sessionId string, t *testing.T) (*response.Writer, *http.Request, database.Connection) {
	w := response.NewWriter(httptest.NewRecorder())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	conn := mocks.MakeMockConnection(t)
	if name != "" {
		if sessionId == "" {
			sessionId = mocks.AdminSesh.ID
		}
		cookie := &http.Cookie{Name: name, Value: sessionId}
		req.AddCookie(cookie)
	}
	return w, req, conn
}

var sessionCookie string = string(cookies.SESSION)

func TestFromCookie(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w, req, conn := setup(sessionCookie, "", t)

		req, ok := FromCookie(w, req, conn)
		assert.Equals(t, ok, true)
		user, err := context.GetUser(req)
		assert.IsNil(t, err)
		assert.DeepEquals(t, user, mocks.Admin)
	})

	t.Run("InvalidCookie", func(t *testing.T) {
		w, req, conn := setup("raisin", "", t)
		reqCopy := *req

		req, ok := FromCookie(w, req, conn)
		assert.Equals(t, ok, false)
		assert.DeepEquals(t, *req, reqCopy)
		assert.Equals(t, w.Status, http.StatusUnauthorized)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		user := mocks.MakeUser()
		sesh := mocks.MakeSession(user.ID)
		w, req, conn := setup(sessionCookie, sesh.ID, t)
		reqCopy := *req
		conn.SetSession(sesh)

		req, ok := FromCookie(w, req, conn)
		assert.Equals(t, ok, false)
		assert.DeepEquals(t, *req, reqCopy)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "user not found during authentication")
	})

	t.Run("CannotSetNewContext", func(t *testing.T) {
		w, req, conn := setup(sessionCookie, "", t)
		buf := logger.MockFileLogger(t)

		req, ok := FromCookie(w, req, conn) // adds user to request context
		reqCopy := *req
		assert.Equals(t, ok, true)
		req, ok = FromCookie(w, req, conn) // request context already has user
		assert.Equals(t, ok, false)
		assert.DeepEquals(t, *req, reqCopy)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		xMessage := "user already in request context"
		assert.Equals(t, w.Body, xMessage)
		assert.Contains(t, buf.String(), fmt.Sprint("[ERROR] ", xMessage))
	})

}

func TestGetUserIdFromCookie(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		userID, err := getUserIDFromCookie(setup(sessionCookie, "", t))
		assert.IsNil(t, err)
		assert.Equals(t, userID, mocks.Admin.ID)
	})

	t.Run("WrongCookie", func(t *testing.T) {
		userID, err := getUserIDFromCookie(setup("raisin", "", t))
		assert.ErrorHasMessage(t, err, "no cookie found with the name beango-session")
		assert.Equals(t, userID, "")
	})

	t.Run("NoCookie", func(t *testing.T) {
		userID, err := getUserIDFromCookie(setup("raisin", "", t))
		assert.ErrorHasMessage(t, err, "no cookie found with the name beango-session")
		assert.Equals(t, userID, "")
	})

	t.Run("NoSession", func(t *testing.T) {
		noSeshUser := mocks.MakeUser()
		buf := logger.MockFileLogger(t)
		w, req, conn := setup(sessionCookie, noSeshUser.ID, t)

		userID, err := getUserIDFromCookie(w, req, conn)
		assert.ErrorHasMessage(t, err, "cookie or session is invalid")
		assert.Equals(t, userID, "")
		assert.Equals(t, buf.String(), "")
		resCookie := w.Header().Get("Set-Cookie")
		xResCookie := fmt.Sprintf(
			"%s=%s; Path=/; Expires=%s; HttpOnly; Secure; SameSite=Strict",
			string(cookies.SESSION),
			"",
			time.Unix(0, 0).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"),
		)
		assert.Equals(t, resCookie, xResCookie)
	})
}
