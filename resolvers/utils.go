package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
	"github.com/raphael-p/beango/utils/validate"
)

// Decodes JSON from HTTP request body and binds it to a struct pointer.
// Writes an HTTP error response on failure.
func getRequestBody(w *response.Writer, r *http.Request, ptr any) bool {
	value := reflect.ValueOf(ptr)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		errorResponse := fmt.Sprintf(
			"expected `ptr` to be a pointer to a struct, got %T",
			ptr,
		)
		w.WriteString(http.StatusBadRequest, errorResponse)
		return false
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(ptr); err != nil {
		errorResponse := fmt.Sprint("malformed request body: ", err)
		w.WriteString(http.StatusBadRequest, errorResponse)
		return false
	}

	fields, err := validate.StructFromJSON(ptr)
	if err != nil {
		w.WriteString(http.StatusBadRequest, err.Error())
		return false
	}
	if len(fields) != 0 {
		errorResponse := fmt.Sprintf("missing required field(s): %s", fields)
		w.WriteString(http.StatusBadRequest, errorResponse)
		return false
	}
	return true
}

// Gets all requested context attached to a request.
// Writes an HTTP error response + logs on failure.
func getRequestContext(
	w *response.Writer,
	r *http.Request,
	keys ...string,
) (*database.User, map[string]string, bool) {
	user, err := context.GetUser(r)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, "failed to fetch request user")
		return nil, nil, false
	}

	params := make(map[string]string)
	for _, key := range keys {
		value, err := context.GetParam(r, key)
		if err != nil {
			logger.Error(err.Error())
			w.WriteString(
				http.StatusInternalServerError,
				fmt.Sprint("failed to fetch path parameter: ", key),
			)
			return nil, nil, false
		}
		params[key] = value
	}

	return user, params, true
}

// Calls `getRequestBody()` then, if successful, `getRequestContext()`
func getRequestBodyAndContext(
	w *response.Writer,
	r *http.Request,
	ptr any,
	keys ...string,
) (*database.User, map[string]string, bool) {
	if ok := getRequestBody(w, r, ptr); !ok {
		return nil, nil, false
	}
	return getRequestContext(w, r, keys...)
}

// Handles an unexpected error from the database
func HandleDatabaseError(w *response.Writer, err error) {
	if err == nil {
		return
	}
	message := "database operation failed"
	logger.Error(message + ": " + err.Error())
	w.WriteString(http.StatusInternalServerError, message)
}
