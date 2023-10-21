package routing

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/server/authenticate"
	"github.com/raphael-p/beango/utils/response"
)

type Middleware func(w *response.Writer, r *http.Request, conn database.Connection) (*http.Request, bool)

var Auth Middleware = func(w *response.Writer, newRequest *http.Request, conn database.Connection) (*http.Request, bool) {
	newRequest, httpError := authenticate.Auth(w, newRequest, conn)
	return newRequest, !resolverutils.ProcessHTTPError(w, httpError)
}

var AuthRedirect Middleware = func(w *response.Writer, r *http.Request, conn database.Connection) (*http.Request, bool) {
	newRequest, httpError := authenticate.Auth(w, r, conn)
	if httpError != nil {
		if httpError.Status >= http.StatusInternalServerError {
			// TODO: better handling for this
			w.WriteString(httpError.Status, httpError.Message)
		} else if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/login")
		} else {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusSeeOther)
		}
		return newRequest, false
	}

	return newRequest, true
}
