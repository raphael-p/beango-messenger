package routing

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/server/authenticate"
	"github.com/raphael-p/beango/utils/response"
)

type Middleware func(w *response.Writer, r *http.Request, conn database.Connection) (*http.Request, bool)

var Auth Middleware = authenticate.FromCookie

var AuthRedirect Middleware = func(w *response.Writer, r *http.Request, conn database.Connection) (*http.Request, bool) {
	if newRequest, ok := authenticate.FromCookie(w, r, conn); ok {
		return newRequest, true
	} else {
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/home")
		} else {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusSeeOther)
		}
		return newRequest, false
	}
}
