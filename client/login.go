package client

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/path"
	"github.com/raphael-p/beango/utils/response"
)

func Login(w *response.Writer, r *http.Request, conn database.Connection) {
	// TODO: redirect to home page if session cookie is present
	// TODO: make single page website https://stackoverflow.com/questions/72914578/single-page-application-with-htmx-url-browsing-history-and-manual-reloading-of
	path, ok := path.RelativeJoin("login.html")
	if !ok {
		message := "could not resolve filepath"
		logger.Error(message)
		w.WriteString(http.StatusInternalServerError, message)
		return
	}
	http.ServeFile(w, r, path)

}

func SubmitLogin(w *response.Writer, r *http.Request, conn database.Connection) {
	action, _ := context.GetParam(r, "action")

	var requests []*http.Request
	if action == "signup" {
		if req, ok := cloneRequest(w, r, 2); ok {
			requests = req
		} else {
			return
		}
	} else {
		requests = append(requests, r)
	}

	if action == "signup" {
		resolvers.CreateUser(w, requests[1], conn)
		if w.Status != http.StatusCreated {
			displayError(w, w.Body)
			return
		}
	}

	resolvers.CreateSession(w, requests[0], conn)
	if w.Status == http.StatusNoContent {
		w.Header().Set("HX-Redirect", "/test")
	}
	displayError(w, w.Body)
}
