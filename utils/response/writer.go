package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Writer struct {
	Status       int
	HideResponse bool
	Body         string
	Time         int64
	http.ResponseWriter
}

func NewWriter(w http.ResponseWriter) *Writer {
	return &Writer{ResponseWriter: w}
}

func (w *Writer) WriteHeader(code int) {
	w.Status = code
}

func (w *Writer) writeBody(body string) {
	w.Body = body
}

func (w *Writer) WriteString(code int, response string) {
	w.WriteHeader(code)
	w.writeBody(response)
}

func (w *Writer) Write(bytes []byte) (int, error) {
	w.HideResponse = true
	w.Body = string(bytes)
	w.Header().Set(
		"Content-Length",
		strconv.FormatInt(int64(len(bytes)), 10),
	)
	return len(bytes), nil
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
	if w.Body != "" && !w.HideResponse {
		out = fmt.Sprintf("%s\n\tresponse: %s", out, w.Body)
	}
	return out
}
