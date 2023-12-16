package resolvers

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/response"
)

func Logout(w *response.Writer, r *http.Request, conn database.Connection) {
	sessionID, _ := cookies.Get(r, cookies.SESSION)
	if sessionID != "" {
		conn.DeleteSession(sessionID)
	}

	cookies.Invalidate(w, cookies.SESSION)
	w.Redirect("/login", r)
}
