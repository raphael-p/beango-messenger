package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Writer struct {
	Status int
	Body   string
	Time   int64
	http.ResponseWriter
}

func NewWriter(w http.ResponseWriter) *Writer {
	return &Writer{ResponseWriter: w}
}

func (w *Writer) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *Writer) writeBody(body []byte) (int, error) {
	w.Body = string(body)
	return w.ResponseWriter.Write(body)
}

func (w *Writer) WriteString(code int, response string) {
	w.WriteHeader(code)
	w.writeBody([]byte(response))
}

func (w *Writer) WriteJSON(code int, responseObject any) {
	w.WriteHeader(code)
	response, err := json.Marshal(responseObject)
	if err != nil {
		w.WriteString(http.StatusBadRequest, err.Error())
	}
	w.Header().Set("content-type", "application/json")
	w.writeBody(response)
}

func (w *Writer) String() string {
	out := fmt.Sprintf("status %d (took %dms)", w.Status, w.Time)
	if w.Body != "" {
		out = fmt.Sprintf("%s\n\tresponse: %s", out, w.Body)
	}
	return out
}
