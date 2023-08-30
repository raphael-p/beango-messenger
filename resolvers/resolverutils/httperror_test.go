package resolverutils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func TestProcessHTTPError(t *testing.T) {
	t.Run("WithError", func(t *testing.T) {
		xError := &HTTPError{100, "this is a message"}
		w := response.NewWriter(httptest.NewRecorder())

		hasError := ProcessHTTPError(w, xError)
		assert.Equals(t, hasError, true)
		assert.Equals(t, w.Status, xError.Status)
		assert.Equals(t, string(w.Body), xError.Message)
	})

	t.Run("WithoutError", func(t *testing.T) {
		w := response.NewWriter(httptest.NewRecorder())

		hasError := ProcessHTTPError(w, nil)
		assert.Equals(t, hasError, false)
		assert.Equals(t, w.Status, 0)
		assert.Equals(t, string(w.Body), "")
	})
}

func TestHandleDatabaseError(t *testing.T) {
	t.Run("WithError", func(t *testing.T) {
		buf := logger.MockFileLogger(t)
		errPrefix := "database operation failed"
		errMessage := "this did not go well"

		httpError := HandleDatabaseError(errors.New(errMessage))
		assert.Equals(t, httpError.Status, http.StatusInternalServerError)
		assert.Equals(t, httpError.Message, errPrefix)
		assert.Contains(t, buf.String(), "[ERROR] "+errPrefix+": "+errMessage)
	})

	t.Run("WithoutError", func(t *testing.T) {
		buf := logger.MockFileLogger(t)

		assert.IsNil(t, HandleDatabaseError(nil))
		assert.Equals(t, buf.String(), "")
	})
}
