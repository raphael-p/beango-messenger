package resolvers

import (
	"html/template"
	"net/http"

	"github.com/raphael-p/beango/client"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func Login(w *response.Writer, r *http.Request, conn database.Connection) {
	if sessionID, err := cookies.Get(r, cookies.SESSION); err == nil {
		if _, ok := conn.CheckSession(sessionID); ok {
			w.Header().Set("Location", "/home")
			w.WriteHeader(http.StatusSeeOther)
			return
		}
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Write([]byte("<div id='content' hx-swap-oob='innerHTML'>" + client.LoginPage + "</div>"))
		return
	}

	skeleton, err := client.GetSkeleton()
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return
	}

	data := map[string]any{"content": template.HTML(client.LoginPage)}
	if err := skeleton.Execute(w, data); err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
	}
}

func SubmitLogin(w *response.Writer, r *http.Request, conn database.Connection) {
	action, _ := context.GetParam(r, "action")
	var input CreateUserInput
	if resolverutils.DisplayHTTPError(w, resolverutils.GetRequestBody(r, &input)) {
		return
	}

	if action == "signup" {
		_, httpError := createUserDatabase(input.Username, input.DisplayName.Value, input.Password, conn)
		if resolverutils.DisplayHTTPError(w, httpError) {
			return
		}
	}

	userID, httpError := checkCredentials(input.Username, input.Password, conn)
	if resolverutils.DisplayHTTPError(w, httpError) {
		return
	}

	if resolverutils.DisplayHTTPError(w, setSession(w, makeSession(userID), conn)) {
		return
	}

	w.Header().Set("HX-Redirect", "/home")
	w.WriteHeader(http.StatusNoContent)
}
