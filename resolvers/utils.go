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
func bindRequestJSON(w *response.Writer, r *http.Request, ptr any) bool {
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
// `ptr` must be a pointer to a struct where all fields are strings.
// Writes an HTTP error response + logs on failure.
func getRequestContext(w *response.Writer, r *http.Request, ptr any) (*database.User, bool) {
	user, err := context.GetUser(r)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, "failed to fetch user")
		return nil, false
	}

	if !validate.PointerToStringStruct(ptr) {
		logger.Error("path param variable must point to a struct of strings")
		w.WriteString(http.StatusInternalServerError, "failed to fetch path parameters")
		return nil, false
	}

	reflectValue := reflect.ValueOf(ptr).Elem()
	reflectType := reflectValue.Type()
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectValue.Field(i)
		key := reflectType.Field(i).Name
		value, err := context.GetParam(r, key)
		if err != nil {
			logger.Error(err.Error())
			w.WriteString(
				http.StatusInternalServerError,
				fmt.Sprint("failed to fetch path parameter: ", key),
			)
			return nil, false
		}
		field.SetString(value)
	}
	return user, true
}
