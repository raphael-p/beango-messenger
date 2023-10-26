package resolvers

import (
	"net/http"

	"github.com/raphael-p/beango/client"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
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

	client.ServeTemplate(w, "loginPage", client.Skeleton+client.LoginPage, nil)
}

func SubmitLogin(w *response.Writer, r *http.Request, conn database.Connection) {
	// getting route param directly instead of using GetRequestContext()
	// because it would error since no user is in the context
	action, _ := context.GetParam(r, resolverutils.ACTION_KEY)
	var input createUserInput
	if resolverutils.DisplayHTTPError(w, resolverutils.GetRequestBody(r, &input)) {
		return
	}
	if resolverutils.DisplayHTTPError(w, createUserInputValidation(&input)) {
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
