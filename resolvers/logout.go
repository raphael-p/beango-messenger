package resolvers

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/response"
)

// TODO: test
func Logout(w *response.Writer, r *http.Request, conn database.Connection) {
	sessionID, _ := cookies.Get(r, cookies.SESSION)
	if sessionID != "" {
		conn.DeleteSession(sessionID)
	}

	cookies.Invalidate(w, cookies.SESSION)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
	} else {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusSeeOther)
	}
}
