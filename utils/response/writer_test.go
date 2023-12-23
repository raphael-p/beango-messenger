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
		w := NewWriter(recorder)

		xStatus := http.StatusAccepted
		w.WriteHeader(xStatus)
		assert.Equals(t, w.Status, xStatus)
		assert.Equals(t, recorder.Code, xStatus)
	})
}

func TestWriteString(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		w := NewWriter(recorder)

		xStatus := http.StatusAccepted
		xBody := "Hello, world!"
		w.WriteString(xStatus, xBody)
		assert.Equals(t, w.Status, xStatus)
		assert.Equals(t, string(w.Body), xBody)
		assert.Equals(t, recorder.Code, xStatus)
		assert.Equals(t, recorder.Body.String(), xBody)
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w := NewWriter(httptest.NewRecorder())
		xStatus := http.StatusAccepted

		w.WriteJSON(xStatus, map[string]string{"message": "Hello, world!"})
		assert.Equals(t, w.Header().Get("Content-Type"), "application/json")
		assert.Equals(t, w.Status, xStatus)
		assert.Equals(t, string(w.Body), "{\"message\":\"Hello, world!\"}")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		w := NewWriter(httptest.NewRecorder())
		body := func(string) { fmt.Println("hello, world!") }

		w.WriteJSON(http.StatusAccepted, body)
		xStatus := http.StatusBadRequest
		assert.Equals(t, w.Status, xStatus)
		assert.Equals(t, string(w.Body), "json: unsupported type: func(string)")
	})
}

func TestWriteHTML(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w := NewWriter(httptest.NewRecorder())
		xStatus := http.StatusAccepted
		xBody := "<div class: >just some text, doesn't have to be valid HTML"

		w.WriteHTML(xStatus, xBody)
		assert.Equals(t, w.Header().Get("Content-Type"), "text/html")
		assert.Equals(t, w.Status, xStatus)
		assert.Equals(t, string(w.Body), xBody)
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
		w := NewWriter(recorder)

		xEvent := "event-one"
		xData := "the event's data"
		xBody := appendToBody("", xEvent, xData)

		err := w.WriteSSE(xEvent, xData)
		assert.IsNil(t, err)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Equals(t, string(w.Body), xBody)
		assert.Equals(t, recorder.Flushed, true)
	})

	t.Run("WriteTwice", func(t *testing.T) {
		w := NewWriter(httptest.NewRecorder())

		xEvent1 := "event-one"
		xData1 := "the event's data"
		xBody := appendToBody("", xEvent1, xData1)

		xEvent2 := "event-two"
		xData2 := "another payload!"
		xBody = appendToBody(xBody, xEvent2, xData2)

		w.WriteSSE(xEvent1, xData1)
		w.WriteSSE(xEvent2, xData2)
		assert.Equals(t, string(w.Body), xBody)
		assert.Equals(t, string(w.Body), xBody)
	})

	t.Run("DefaultEvent", func(t *testing.T) {
		w := NewWriter(httptest.NewRecorder())
		xEvent := "message"
		xData := "the event's data"
		xBody := appendToBody("", xEvent, xData)

		w.WriteSSE("", xData)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Equals(t, string(w.Body), xBody)
	})

	t.Run("NoFlushResponseWriter", func(t *testing.T) {
		w := NewWriter(noFlushRecorder{httptest.NewRecorder()})

		err := w.WriteSSE("", "")
		assert.ErrorHasMessage(t, err, "response writer does not have a flusher")
	})
}

func TestRedirect(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := NewWriter(httptest.NewRecorder())
		xLocation := "/somewhere/else"
		xStatus := http.StatusSeeOther

		w.Redirect(xLocation, r)
		assert.Equals(t, w.Header().Get("Location"), xLocation)
		assert.Equals(t, w.Status, xStatus)
	})

	t.Run("FromHTMX", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := NewWriter(httptest.NewRecorder())
		r.Header.Set("HX-Request", "true")
		xLocation := "/somewhere/else"
		xStatus := http.StatusOK

		w.Redirect(xLocation, r)
		assert.Equals(t, w.Header().Get("Location"), "")
		assert.Equals(t, w.Header().Get("HX-Redirect"), xLocation)
		assert.Equals(t, w.Status, xStatus)
	})

}

func TestString(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w := NewWriter(httptest.NewRecorder())
		w.Status = http.StatusOK
		w.Time = 100 * time.Millisecond.Milliseconds()
		w.Body = []byte("Hello, world!")

		xOut := fmt.Sprintf("status %d (took %dms)\n\tresponse: %s", http.StatusOK, w.Time, w.Body)
		assert.Equals(t, w.String(), xOut)
	})

	t.Run("BodyHidden", func(t *testing.T) {
		w := NewWriter(httptest.NewRecorder())
		w.Status = http.StatusOK
		w.Time = 100 * time.Millisecond.Milliseconds()
		w.Body = []byte("Hello, world!")
		w.HideResponse = true

		xOut := fmt.Sprintf("status %d (took %dms)", http.StatusOK, w.Time)
		assert.Equals(t, w.String(), xOut)
	})
}
