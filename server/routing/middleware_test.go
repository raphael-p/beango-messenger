package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/response"
)

func TestAuth(t *testing.T) {
	setup := func(t *testing.T) (*response.Writer, *http.Request, database.Connection) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := response.NewWriter(httptest.NewRecorder())
		conn := mocks.MakeMockConnection()
		return w, req, conn
	}

	t.Run("AuthSucceeds", func(t *testing.T) {
		w, req, conn := setup(t)
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: mocks.AdminSesh.ID}
		req.AddCookie(cookie)

		newReq, proceed := Auth(w, req, conn)
		assert.Equals(t, proceed, true)
		user, err := context.GetUser(newReq)
		assert.IsNil(t, err)
		assert.DeepEquals(t, user, mocks.Admin)
	})

	t.Run("AuthFails", func(t *testing.T) {
		w, req, conn := setup(t)

		_, proceed := Auth(w, req, conn)
		assert.Equals(t, proceed, false)
		assert.Equals(t, w.Status, http.StatusUnauthorized)
	})
}

func TestAuthRedirect(t *testing.T) {
	setup := func(t *testing.T) (*response.Writer, *http.Request, database.Connection) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := response.NewWriter(httptest.NewRecorder())
		conn := mocks.MakeMockConnection()
		return w, req, conn
	}

	t.Run("NoRedirectOnSuccess", func(t *testing.T) {
		w, req, conn := setup(t)
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: mocks.AdminSesh.ID}
		req.AddCookie(cookie)

		_, proceed := AuthRedirect(w, req, conn)
		assert.Equals(t, proceed, true)
		assert.Equals(t, w.Status, 0)
		assert.Equals(t, w.Header().Get("Location"), "")
		assert.Equals(t, w.Header().Get("HX-Redirect"), "")
	})

	t.Run("RedirectOnFailure", func(t *testing.T) {
		w, req, conn := setup(t)

		_, proceed := AuthRedirect(w, req, conn)
		assert.Equals(t, proceed, false)
		assert.Equals(t, w.Status, http.StatusSeeOther)
		assert.Equals(t, w.Header().Get("Location"), "/login")
	})

	t.Run("HXRedirectOnFailure", func(t *testing.T) {
		w, req, conn := setup(t)
		req.Header.Set("HX-Request", "true")

		_, proceed := AuthRedirect(w, req, conn)
		assert.Equals(t, proceed, false)
		assert.Equals(t, w.Status, 0)
		assert.Equals(t, w.Header().Get("HX-Redirect"), "/login")
	})
}
