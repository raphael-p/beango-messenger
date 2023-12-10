package resolvers

import (
	"net/http"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/response"
)

func TestLogout(t *testing.T) {
	checkValidResponse := func(w *response.Writer, r *http.Request, conn database.Connection) {
		Logout(w, r, conn)
		assert.Equals(t, w.Status, http.StatusSeeOther)
		assert.Equals(t, w.Header().Get("Location"), "/login")
		assert.Contains(t, w.Header().Get("Set-Cookie"), "Expires=Thu, 01 Jan 1970")
	}

	t.Run("NoSessionCookie", func(t *testing.T) {
		checkValidResponse(resolverutils.CommonSetup(""))
	})

	t.Run("InvalidSessionCookie", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: "not_a_valid_num"}
		r.AddCookie(cookie)

		checkValidResponse(w, r, conn)
	})

	t.Run("ValidSessionCookie", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		user, _ := conn.SetUser(mocks.MakeUser())
		xSesh := mocks.MakeSession(user.ID)
		conn.SetSession(xSesh)
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: xSesh.ID}
		r.AddCookie(cookie)

		sesh, _ := conn.GetSessionByUserID(xSesh.UserID)
		assert.IsNotNil(t, sesh)

		checkValidResponse(w, r, conn)
		sesh, _ = conn.GetSessionByUserID(xSesh.UserID)
		assert.IsNil(t, sesh)
	})

	t.Run("FromHTMX", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		r.Header.Set("HX-Request", "true")

		Logout(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Equals(t, w.Header().Get("HX-Redirect"), "/login")
		assert.Contains(t, w.Header().Get("Set-Cookie"), "Expires=Thu, 01 Jan 1970")
	})
}
