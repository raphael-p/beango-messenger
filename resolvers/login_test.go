package resolvers

import (
	"net/http"
	"testing"

	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/cookies"
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
