package routing

import (
	"net/http"
	"testing"

	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
)

func TestAuth(t *testing.T) {
	t.Run("AuthSucceeds", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: mocks.AdminSesh.ID}
		req.AddCookie(cookie)

		newReq, proceed := Auth(w, req, conn)
		assert.Equals(t, proceed, true)
		user, err := context.GetUser(newReq)
		assert.IsNil(t, err)
		assert.DeepEquals(t, user, mocks.Admin)
	})

	t.Run("AuthFails", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")

		_, proceed := Auth(w, req, conn)
		assert.Equals(t, proceed, false)
		assert.Equals(t, w.Status, http.StatusUnauthorized)
	})
}

func TestAuthRedirect(t *testing.T) {
	t.Run("NoRedirectOnSuccess", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: mocks.AdminSesh.ID}
		req.AddCookie(cookie)

		_, proceed := AuthRedirect(w, req, conn)
		assert.Equals(t, proceed, true)
		assert.Equals(t, w.Status, 0)
		assert.Equals(t, w.Header().Get("Location"), "")
		assert.Equals(t, w.Header().Get("HX-Redirect"), "")
	})

	t.Run("RedirectOnFailure", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")

		_, proceed := AuthRedirect(w, req, conn)
		assert.Equals(t, proceed, false)
		assert.Equals(t, w.Status, http.StatusSeeOther)
		assert.Equals(t, w.Header().Get("Location"), "/login")
	})

	t.Run("HXRedirectOnFailure", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")
		req.Header.Set("HX-Request", "true")

		_, proceed := AuthRedirect(w, req, conn)
		assert.Equals(t, proceed, false)
		assert.Equals(t, w.Status, 200)
		assert.Equals(t, w.Header().Get("HX-Redirect"), "/login")
	})
}

func TestAuthWeak(t *testing.T) {
	t.Run("AuthSucceeds", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: mocks.AdminSesh.ID}
		req.AddCookie(cookie)

		newReq, proceed := AuthWeak(w, req, conn)
		assert.Equals(t, proceed, true)
		user, err := context.GetUser(newReq)
		assert.IsNil(t, err)
		assert.DeepEquals(t, user, mocks.Admin)
	})

	t.Run("AuthFails", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")

		newReq, proceed := AuthWeak(w, req, conn)
		assert.Equals(t, proceed, true)
		_, err := context.GetUser(newReq)
		assert.ErrorHasMessage(t, err, "user not found in request context")
	})
}
