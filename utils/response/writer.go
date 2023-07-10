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
}

func (w *Writer) writeBody(body string) {
	w.Body = string(body)
}

func (w *Writer) WriteString(code int, response string) {
	w.WriteHeader(code)
	w.writeBody(response)
}

func (w *Writer) WriteJSON(code int, responseObject any) {
	response, err := json.Marshal(responseObject)
	if err != nil {
		errCode := http.StatusBadRequest
		w.WriteHeader(errCode)
		w.WriteString(errCode, err.Error())
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.writeBody(string(response))
}

func (w *Writer) Commit() (int, error) {
	w.ResponseWriter.WriteHeader(w.Status)
	return w.ResponseWriter.Write([]byte(w.Body))
}

func (w *Writer) String() string {
	out := fmt.Sprintf("status %d (took %dms)", w.Status, w.Time)
	if w.Body != "" {
		out = fmt.Sprintf("%s\n\tresponse: %s", out, w.Body)
	}
	return out
}
