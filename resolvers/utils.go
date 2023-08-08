package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
	"github.com/raphael-p/beango/utils/validate"
)

// Decodes JSON from HTTP request body and binds it to a struct pointer.
// Writes an HTTP error response on failure.
func getRequestBody(r *http.Request, ptr any) *HTTPError {
	value := reflect.ValueOf(ptr)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		errorResponse := fmt.Sprintf(
			"expected `ptr` to be a pointer to a struct, got %T",
			ptr,
		)
		return &HTTPError{http.StatusBadRequest, errorResponse}
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(ptr); err != nil {
		errorResponse := fmt.Sprint("malformed request body: ", err)
		return &HTTPError{http.StatusBadRequest, errorResponse}
	}

	fields, err := validate.StructFromJSON(ptr)
	if err != nil {
		return &HTTPError{http.StatusBadRequest, err.Error()}
	}
	if len(fields) != 0 {
		errorResponse := fmt.Sprintf("missing required field(s): %s", fields)
		return &HTTPError{http.StatusBadRequest, errorResponse}
	}
	return nil
}

// Gets all requested context attached to a request.
// Writes an HTTP error response + logs on failure.
func getRequestContext(r *http.Request, paramKeys ...string) (*database.User, *RouteParams, *HTTPError) {
	user, err := context.GetUser(r)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, &HTTPError{http.StatusInternalServerError, "failed to fetch request user"}
	}

	routeParams, httpError := MakeRouteParams(r, paramKeys...)
	if httpError != nil {
		return nil, nil, httpError
	}

	return user, routeParams, nil
}

// Calls `getRequestBody()` then, if successful, `getRequestContext()`
func getRequestBodyAndContext(
	r *http.Request,
	ptr any,
	paramKeys ...string,
) (*database.User, *RouteParams, *HTTPError) {
	if httpError := getRequestBody(r, ptr); httpError != nil {
		return nil, nil, httpError
	}
	return getRequestContext(r, paramKeys...)
}

type HTTPError struct {
	status  int
	message string
}

// Writes message and status of HTTPError to the response
// Returns false if httpError is nil, true otherwise
// TODO: test
func ProcessHTTPError(w *response.Writer, httpError *HTTPError) bool {
	if httpError == nil {
		return false
	}
	w.WriteString(httpError.status, httpError.message)
	return true
}

// Handles an unexpected error from the database
func HandleDatabaseError(err error) *HTTPError {
	if err == nil {
		return nil
	}
	message := "database operation failed"
	logger.Error(message + ": " + err.Error())
	return &HTTPError{http.StatusInternalServerError, message}
}

// TODO: unit test
func StringToInt(str string, bitSize int) (int64, *HTTPError) {
	num, err := strconv.ParseInt("1", 10, bitSize)
	if err != nil {
		return 0, &HTTPError{http.StatusBadRequest, "chat ID must be an integer"}
	}
	return num, nil
}
