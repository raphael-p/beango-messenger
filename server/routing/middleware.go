package routing

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/server/authenticate"
	"github.com/raphael-p/beango/utils/response"
)

type Middleware func(w *response.Writer, r *http.Request, conn database.Connection) (*http.Request, bool)

// Adds user to request context.
// On failure, returns a 401.
var Auth Middleware = func(w *response.Writer, newRequest *http.Request, conn database.Connection) (*http.Request, bool) {
	newRequest, httpError := authenticate.Auth(w, newRequest, conn)
	return newRequest, !resolverutils.ProcessHTTPError(w, httpError)
}

// Adds user to request context.
// On failure, redirects to /login.
var AuthRedirect Middleware = func(w *response.Writer, r *http.Request, conn database.Connection) (*http.Request, bool) {
	newRequest, httpError := authenticate.Auth(w, r, conn)
	if httpError == nil {
		return newRequest, true
	}

	if httpError.Status >= http.StatusInternalServerError {
		if r.Header.Get("HX-Request") == "true" {
			resolverutils.DisplayHTTPError(w, httpError)
		} else {
			w.WriteString(httpError.Status, httpError.Message)
		}
	} else {
		w.Redirect("/login", r)
	}
	return newRequest, false
}

// Adds user to request context.
// On failure, it proceeds anyway.
var AuthWeak Middleware = func(w *response.Writer, newRequest *http.Request, conn database.Connection) (*http.Request, bool) {
	newRequest, _ = authenticate.Auth(w, newRequest, conn)
	return newRequest, true
}
