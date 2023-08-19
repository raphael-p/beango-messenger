package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/utils/response"
)

func TestDisplayError(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w := response.NewWriter(httptest.NewRecorder())
		message := "something went wrong!"

		DisplayError(w, message)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), fmt.Sprintf(">%s<", message))
	})

	t.Run("Empty", func(t *testing.T) {
		w := response.NewWriter(httptest.NewRecorder())

		DisplayError(w, "")
		assert.Equals(t, w.Status, 0)
		assert.Contains(t, string(w.Body), "")
	})
}
