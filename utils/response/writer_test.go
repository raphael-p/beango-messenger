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

func TestWriteHTML(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)

		xStatus := http.StatusAccepted
		xBody := "<div class: >just some text, doesn't have to be valid HTML"
		writer.WriteHTML(xStatus, xBody)
		assert.Equals(t, recorder.Header().Get("Content-Type"), "text/html")
		assert.Equals(t, writer.Status, xStatus)
		assert.Equals(t, string(writer.Body), xBody)
	})
}

type noFlushRecorder struct {
	*httptest.ResponseRecorder
}

func (recorder *noFlushRecorder) Flush() {}

func TestWriteSSE(t *testing.T) {
	appendToBody := func(body, event, data string) string {
		return fmt.Sprintf("%sevent: %s\ndata: %s\n\n", body, event, data)
	}

	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)

		xEvent := "event-one"
		xData := "the event's data"
		xBody := appendToBody("", xEvent, xData)

		err := writer.WriteSSE(xEvent, xData)
		assert.IsNil(t, err)
		assert.Equals(t, writer.Status, http.StatusOK)
		assert.Equals(t, string(writer.Body), xBody)
		assert.Equals(t, recorder.Flushed, true)
	})

	t.Run("WriteTwice", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)

		xEvent1 := "event-one"
		xData1 := "the event's data"
		xBody := appendToBody("", xEvent1, xData1)

		xEvent2 := "event-two"
		xData2 := "another payload!"
		xBody = appendToBody(xBody, xEvent2, xData2)

		writer.WriteSSE(xEvent1, xData1)
		writer.WriteSSE(xEvent2, xData2)
		assert.Equals(t, string(writer.Body), xBody)
		assert.Equals(t, string(writer.Body), xBody)
	})

	t.Run("DefaultEvent", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writer := NewWriter(recorder)
		xEvent := "message"
		xData := "the event's data"
		xBody := appendToBody("", xEvent, xData)

		writer.WriteSSE("", xData)
		assert.Equals(t, writer.Status, http.StatusOK)
		assert.Equals(t, string(writer.Body), xBody)
	})

	t.Run("NoFlushResponseWriter", func(t *testing.T) {
		recorder := noFlushRecorder{httptest.NewRecorder()}
		writer := NewWriter(recorder)

		err := writer.WriteSSE("", "")
		assert.ErrorHasMessage(t, err, "response writer does not have a flusher")
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
