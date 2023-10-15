package resolverutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/validate"
)

// Decodes JSON from HTTP request body and binds it to a struct pointer.
// Writes an HTTP error response on failure.
func GetRequestBody(r *http.Request, ptr any) *HTTPError {
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
func GetRequestContext(r *http.Request, paramKeys ...string) (*database.User, *RouteParams, *HTTPError) {
	user, err := context.GetUser(r)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, &HTTPError{http.StatusInternalServerError, "failed to fetch request user"}
	}

	routeParams, httpError := extractRouteParams(r, paramKeys...)
	if httpError != nil {
		return nil, nil, httpError
	}

	return user, routeParams, nil
}

// Calls `resolverutils.GetRequestBody()` then, if successful, `resolverutils.GetRequestContext()`
func GetRequestBodyAndContext(
	r *http.Request,
	ptr any,
	paramKeys ...string,
) (*database.User, *RouteParams, *HTTPError) {
	if httpError := GetRequestBody(r, ptr); httpError != nil {
		return nil, nil, httpError
	}
	return GetRequestContext(r, paramKeys...)
}

// Extracts query parameter from request
// errors if missing when isRequired is true
// errors if the parameter value is an empty string
func GetRequestQueryParam(r *http.Request, key string, isRequired bool) (string, *HTTPError) {
	query := r.URL.Query()
	if !query.Has(key) {
		var httpError *HTTPError
		if isRequired {
			httpError = &HTTPError{http.StatusBadRequest, "missing required query parameter: " + key}
		}
		return "", httpError
	}

	value := query.Get(key)
	if value == "" {
		return "", &HTTPError{http.StatusBadRequest, "query parameter cannot be empty: " + key}
	}
	return value, nil
}

// Calls GetRequestQueryParam, then casts the result to an integer
func GetRequestQueryParamInt(r *http.Request, key string, isRequired bool) (int64, *HTTPError) {
	value, httpError := GetRequestQueryParam(r, key, isRequired)
	if httpError != nil {
		return 0, httpError
	}
	var intValue int64
	var err error
	intValue, err = strconv.ParseInt(value, 10, 64)
	if err != nil {
		message := fmt.Sprintf("query parameter '%s' must be an integer", key)
		return 0, &HTTPError{http.StatusBadRequest, message}
	}
	return intValue, nil
}
