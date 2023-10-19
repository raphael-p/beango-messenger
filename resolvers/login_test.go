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
	"github.com/raphael-p/beango/utils/response"
)

func TestLogin(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")
		req = resolverutils.SetContext(t, req, mocks.Admin, nil)

		Login(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), "<html>", "</html")
	})

	t.Run("ValidSessionCookie", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: mocks.AdminSesh.ID}
		req.AddCookie(cookie)

		Login(w, req, conn)
		assert.Equals(t, w.Status, http.StatusSeeOther)
		assert.Equals(t, string(w.Body), "")
		assert.Equals(t, w.Header().Get("Location"), "/home")
	})

	t.Run("InvalidSessionCookie", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: "not-a-valid-session-id"}
		req.AddCookie(cookie)

		Login(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), "<html>", "</html")
	})
}

func TestSubmitLogin(t *testing.T) {
	config.CreateConfig()

	body := func(username, password string) string {
		return fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	}

	checkSuccessfulLogin := func(w *response.Writer, req *http.Request, conn database.Connection) {
		SubmitLogin(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNoContent)
		assert.Equals(t, string(w.Body), "")
		assert.Equals(t, w.Header().Get("HX-Redirect"), "/home")
	}

	t.Run("NormalWithLogin", func(t *testing.T) {
		mockUser := mocks.MakeUser()
		w, req, conn := resolverutils.CommonSetup(body(mockUser.Username, mocks.PASSWORD))
		conn.SetUser(mockUser)
		param := map[string]string{"action": "login"}
		req = resolverutils.SetContext(t, req, mocks.Admin, param)

		checkSuccessfulLogin(w, req, conn)
	})

	t.Run("NormalWithSignup", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup(body("someNewUser", "123"))
		param := map[string]string{"action": "signup"}
		req = resolverutils.SetContext(t, req, mocks.Admin, param)

		checkSuccessfulLogin(w, req, conn)
	})
}
