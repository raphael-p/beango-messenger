package response

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raphael-p/beango/test/assert"
)

func TestWriteHeader(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)

		xStatus := http.StatusAccepted
		writer.WriteHeader(xStatus)
		assert.Equals(t, writer.Status, xStatus)
		assert.Equals(t, recorder.Code, http.StatusOK) // default
	})
}

func TestWriteString(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)

		xStatus := http.StatusAccepted
		xBody := "Hello, world!"
		writer.WriteString(xStatus, xBody)
		assert.Equals(t, writer.Status, xStatus)
		assert.Equals(t, writer.Body, xBody)
		assert.Equals(t, recorder.Code, http.StatusOK) // default
		assert.Equals(t, recorder.Body.String(), "")   // default
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)

		xStatus := http.StatusAccepted
		writer.WriteJSON(xStatus, map[string]string{"message": "Hello, world!"})
		assert.Equals(t, recorder.Header().Get("content-type"), "application/json")
		assert.Equals(t, writer.Status, xStatus)
		assert.Equals(t, writer.Body, "{\"message\":\"Hello, world!\"}")
		assert.Equals(t, recorder.Code, http.StatusOK) // default
		assert.Equals(t, recorder.Body.String(), "")   // default
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)

		body := func(string) { fmt.Println("hello, world!") }
		writer.WriteJSON(http.StatusAccepted, body)
		xStatus := http.StatusBadRequest
		assert.Equals(t, writer.Status, xStatus)
		assert.Equals(t, writer.Body, "json: unsupported type: func(string)")
		assert.Equals(t, recorder.Code, http.StatusOK) // default
		assert.Equals(t, recorder.Body.String(), "")   // default
	})
}

func TestString(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)
		writer.Status = http.StatusOK
		writer.Time = 100 * time.Millisecond.Milliseconds()
		writer.Body = "Hello, world!"

		xOut := fmt.Sprintf("status %d (took %dms)\n\tresponse: %s", http.StatusOK, writer.Time, writer.Body)
		assert.Equals(t, writer.String(), xOut)
	})
}

func TestCommit(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)
		xStatus := 573
		xBody := "Hello, World!"
		writer.Status = xStatus
		writer.Body = xBody

		assert.Equals(t, recorder.Code, http.StatusOK)
		assert.Equals(t, recorder.Body.String(), "")
		writer.Commit()
		assert.Equals(t, recorder.Code, xStatus)
		assert.Equals(t, recorder.Body.String(), xBody)
	})
}
