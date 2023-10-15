package resolverutils

import (
	"errors"
	"fmt"
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
		AssertHTTPError(t, httpError, http.StatusInternalServerError, errPrefix)
		assert.Contains(t, buf.String(), "[ERROR] "+errPrefix+": "+errMessage)
	})

	t.Run("WithoutError", func(t *testing.T) {
		buf := logger.MockFileLogger(t)

		assert.IsNil(t, HandleDatabaseError(nil))
		assert.Equals(t, buf.String(), "")
	})
}

func TestDisplayHTTPError(t *testing.T) {
	t.Run("WithError", func(t *testing.T) {
		xError := &HTTPError{100, "this is a message"}
		w := response.NewWriter(httptest.NewRecorder())

		hasError := DisplayHTTPError(w, xError)
		assert.Equals(t, hasError, true)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), fmt.Sprintf(">%s<", xError.Message))
	})

	t.Run("WithoutError", func(t *testing.T) {
		w := response.NewWriter(httptest.NewRecorder())

		hasError := DisplayHTTPError(w, nil)
		assert.Equals(t, hasError, false)
		assert.Equals(t, w.Status, 0)
		assert.Equals(t, string(w.Body), "")
	})
}
