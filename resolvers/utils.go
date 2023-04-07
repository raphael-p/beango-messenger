package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
	"github.com/raphael-p/beango/utils/validate"
)

// Decodes JSON from HTTP request body and binds it to a struct pointer.
// Writes an HTTP error response on failure.
func bindRequestJSON(w *response.Writer, r *http.Request, ptr any) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(ptr); err != nil {
		response := fmt.Sprint("malformed request body: ", err)
		w.WriteString(http.StatusBadRequest, response)
		return false
	}
	if fields := validate.DeserialisedJSON(ptr); len(fields) != 0 {
		response := fmt.Sprintf("missing required field(s): %s", fields)
		w.WriteString(http.StatusBadRequest, response)
		return false
	}
	return true
}

// Gets all requested context attached to a request. Only returns user if one of the keys
// is "user". All keys are assumed to be context parameters and are returned in a map.
// Writes an HTTP error response + logs on failure.
func getRequestContext(
	w *response.Writer,
	r *http.Request,
	keys ...string,
) (*database.User, map[string]string, bool) {
	user, err := context.GetUser(r)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, "failed to get user from context")
		return nil, nil, false
	}

	params := make(map[string]string)
	for _, key := range keys {
		value, err := context.GetParam(r, key)
		if err != nil {
			logger.Error(err.Error())
			w.WriteString(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to get %s from context", key),
			)
			return nil, nil, false
		}
		params[key] = value
	}
	return user, params, true
}
