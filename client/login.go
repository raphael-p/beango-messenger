package client

import (
	"html/template"
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func Login(w *response.Writer, r *http.Request, conn database.Connection) {
	// TODO: use htmx to replace if request comes from htmx
	// TODO: redirect to home page if session cookie is present
	if r.Header.Get("HX-Request") == "true" {
		w.Write([]byte("<div id='content' hx-swap-oob='innerHTML'>" + loginPage + "</div>"))
		return
	}

	skeleton, err := getSkeleton()
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}

	data := map[string]any{"content": template.HTML(loginPage)}
	if err := skeleton.Execute(w, data); err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
	}
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
			displayError(w, string(w.Body))
			return
		}
	}

	resolvers.CreateSession(w, requests[0], conn)
	if w.Status == http.StatusNoContent {
		w.Status = 0
		w.Header().Set("HX-Push", "/home")
		Home(w, r, conn)
		return
	}
	displayError(w, string(w.Body))
}
