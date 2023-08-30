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
		assert.Equals(t, recorder.Code, xStatus)
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
		assert.Equals(t, string(writer.Body), xBody)
		assert.Equals(t, recorder.Code, xStatus)
		assert.Equals(t, recorder.Body.String(), xBody)
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)

		xStatus := http.StatusAccepted
		writer.WriteJSON(xStatus, map[string]string{"message": "Hello, world!"})
		assert.Equals(t, recorder.Header().Get("Content-Type"), "application/json")
		assert.Equals(t, writer.Status, xStatus)
		assert.Equals(t, string(writer.Body), "{\"message\":\"Hello, world!\"}")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)

		body := func(string) { fmt.Println("hello, world!") }
		writer.WriteJSON(http.StatusAccepted, body)
		xStatus := http.StatusBadRequest
		assert.Equals(t, writer.Status, xStatus)
		assert.Equals(t, string(writer.Body), "json: unsupported type: func(string)")
	})
}

func TestString(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)
		writer.Status = http.StatusOK
		writer.Time = 100 * time.Millisecond.Milliseconds()
		writer.Body = []byte("Hello, world!")

		xOut := fmt.Sprintf("status %d (took %dms)\n\tresponse: %s", http.StatusOK, writer.Time, writer.Body)
		assert.Equals(t, writer.String(), xOut)
	})

	t.Run("BodyHidden", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)
		writer.Status = http.StatusOK
		writer.Time = 100 * time.Millisecond.Milliseconds()
		writer.Body = []byte("Hello, world!")
		writer.HideResponse = true

		xOut := fmt.Sprintf("status %d (took %dms)", http.StatusOK, writer.Time)
		assert.Equals(t, writer.String(), xOut)
	})
}
