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
	loginPage := template.HTML(`<span class="title"><span>> Beango Messenger </span></span>
	<div id="login-form">
		<form hx-ext='json-enc'>
			<div class="form-row">
				<label for="username">Username:</label>
				<input type="text" name="username">
			</div>
			<div class="form-row">
				<label for="password">Password:</label>
				<input type="password" name="password">
			</div>
			<div class="form-row button-row">
				<button hx-post="/login/login" type="submit" hx-swap="none">Log In</button>
				<button hx-post="/login/signup" type="submit" hx-swap="none">Sign Up</button>
			</div>
			<div id="errors" class="error"></div>
		</form>
	</div>`)

	if err := container.Execute(w, loginPage); err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusOK)
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
		w.Header().Set("HX-Redirect", "/test")
	}
	displayError(w, string(w.Body))
}
