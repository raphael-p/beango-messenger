package resolvers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func TestCreateSession(t *testing.T) {
	config.CreateConfig()
	setup := func(username, password string) (
		*response.Writer,
		*http.Request,
		database.Connection,
	) {
		body := fmt.Sprintf(`{"Username": "%s", "Password": "%s"}`, username, password)
		w, req := resolverutils.MockRequest(body)
		return w, req, mocks.MakeMockConnection()
	}

	t.Run("Normal", func(t *testing.T) {
		w, req, conn := setup(mocks.ADMIN_USERNAME, mocks.PASSWORD)

		CreateSession(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNoContent)
		assert.HasLength(t, w.Header()["Set-Cookie"], 1)
	})

	t.Run("RequestHasInvalidSession", func(t *testing.T) {
		w, req, conn := setup(mocks.ADMIN_USERNAME, mocks.PASSWORD)
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: mocks.AdminSesh.ID}
		req.AddCookie(cookie)
		conn.DeleteSession(mocks.AdminSesh.ID)

		CreateSession(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNoContent)
		assert.HasLength(t, w.Header()["Set-Cookie"], 1)
	})

	t.Run("RequestHasValidSession", func(t *testing.T) {
		w, req, conn := setup(mocks.ADMIN_USERNAME, mocks.PASSWORD)
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: mocks.AdminSesh.ID}
		req.AddCookie(cookie)

		CreateSession(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNoContent)
		assert.Equals(t, string(w.Body), "")
	})

	t.Run("WrongUsername", func(t *testing.T) {
		w, req, conn := setup(mocks.ADMIN_USERNAME+" ", mocks.PASSWORD)

		CreateSession(w, req, conn)
		assert.Equals(t, w.Status, http.StatusUnauthorized)
		assert.Equals(t, string(w.Body), "login credentials are incorrect")
		assert.HasLength(t, w.Header()["Set-Cookie"], 0)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		w, req, conn := setup(mocks.ADMIN_USERNAME, mocks.PASSWORD+" ")

		CreateSession(w, req, conn)
		assert.Equals(t, w.Status, http.StatusUnauthorized)
		assert.Equals(t, string(w.Body), "login credentials are incorrect")
		assert.HasLength(t, w.Header()["Set-Cookie"], 0)
	})

	t.Run("ResponseAlreadySetsCookie", func(t *testing.T) {
		w, req, conn := setup(mocks.ADMIN_USERNAME, mocks.PASSWORD)
		http.SetCookie(w, &http.Cookie{Name: string(cookies.SESSION)})
		buf := logger.MockFileLogger(t)

		CreateSession(w, req, conn)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, string(w.Body), "failed to create session cookie")
		xError := fmt.Sprint(
			"failed to create session cookie: response header already ",
			"sets a cookie with the name ",
			string(cookies.SESSION),
		)
		assert.Contains(t, buf.String(), "[ERROR]", xError)
	})
}
