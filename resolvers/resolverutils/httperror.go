package resolverutils

import (
	"net/http"

	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

type HTTPError struct {
	Status  int
	Message string
}

// Writes message and status of HTTPError to the response
// Returns false if httpError is nil, true otherwise
func ProcessHTTPError(w *response.Writer, httpError *HTTPError) bool {
	if httpError == nil {
		return false
	}
	w.WriteString(httpError.Status, httpError.Message)
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
