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

func TestCheckCredentials(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		userID, httpError := checkCredentials(mocks.ADMIN_USERNAME, mocks.PASSWORD, conn)
		assert.IsNil(t, httpError)
		assert.Equals(t, userID, mocks.ADMIN_ID)
	})

	t.Run("WrongUsername", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		userID, httpError := checkCredentials(mocks.ADMIN_USERNAME+" ", mocks.PASSWORD, conn)
		assert.Equals(t, userID, 0)
		assert.IsNotNil(t, httpError)
		assert.Equals(t, httpError.Status, http.StatusUnauthorized)
		assert.Equals(t, httpError.Message, "login credentials are incorrect")
	})

	t.Run("WrongPassword", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		userID, httpError := checkCredentials(mocks.ADMIN_USERNAME, mocks.PASSWORD+" ", conn)
		assert.Equals(t, userID, 0)
		assert.IsNotNil(t, httpError)
		assert.Equals(t, httpError.Status, http.StatusUnauthorized)
		assert.Equals(t, httpError.Message, "login credentials are incorrect")
	})
}

func TestSetSession(t *testing.T) {
	setup := func(username, password string) (
		*response.Writer,
		database.Connection,
	) {
		body := fmt.Sprintf(`{"Username": "%s", "Password": "%s"}`, username, password)
		w, _ := resolverutils.MockRequest(body)
		return w, mocks.MakeMockConnection()
	}

	t.Run("Normal", func(t *testing.T) {
		w, conn := setup(mocks.ADMIN_USERNAME, mocks.PASSWORD)

		session := &database.Session{ID: "569"}
		httpError := setSession(w, session, conn)
		assert.IsNil(t, httpError)
		setCookieHeader := w.Header()["Set-Cookie"]
		assert.HasLength(t, setCookieHeader, 1)
		assert.Contains(t, setCookieHeader[0], session.ID)
	})

	t.Run("SessionCookieAlreadySet", func(t *testing.T) {
		w, conn := setup(mocks.ADMIN_USERNAME, mocks.PASSWORD)
		http.SetCookie(w, &http.Cookie{Name: string(cookies.SESSION)})
		buf := logger.MockFileLogger(t)

		httpError := setSession(w, &database.Session{}, conn)
		assert.IsNotNil(t, httpError)
		assert.Equals(t, httpError.Status, http.StatusInternalServerError)
		assert.Equals(t, httpError.Message, "failed to create session cookie")
		xError := fmt.Sprint(
			"failed to create session cookie: response header already ",
			"sets a cookie with the name ",
			string(cookies.SESSION),
		)
		assert.Contains(t, buf.String(), "[ERROR]", xError)
	})
}

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
		assert.HasLength(t, w.Header()["Set-Cookie"], 0)
	})
}
