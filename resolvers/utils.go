package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/httputils"
	"github.com/raphael-p/beango/utils"
	"github.com/raphael-p/beango/validators"
)

// Decodes JSON from HTTP request body and binds it to a struct pointer
func bindRequestJSON(w *httputils.ResponseWriter, r *http.Request, ptr any) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(ptr); err != nil {
		response := fmt.Sprint("malformed request body: ", err)
		w.StringResponse(http.StatusBadRequest, response)
		return false
	}
	if fields := validators.DeserialisedJSON(ptr); len(fields) != 0 {
		response := fmt.Sprintf("missing required field(s): %s", fields)
		w.StringResponse(http.StatusBadRequest, response)
		return false
	}
	return true
}

func extractUser(r *http.Request) (*database.User, error) {
	rawUser, err := httputils.GetUserFromContext(r)
	if err != nil {
		return nil, err
	}
	user, ok := rawUser.(*database.User)
	if !ok {
		message := "context user not of type User"
		utils.Logger.Error(message)
		return nil, fmt.Errorf(message)
	}
	return user, nil
}
