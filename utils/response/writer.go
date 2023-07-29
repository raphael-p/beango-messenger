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
	Body         []byte
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
	w.Body = []byte(body)
}

func (w *Writer) WriteString(code int, response string) {
	w.WriteHeader(code)
	w.writeBody(response)
}

func (w *Writer) Write(bytes []byte) (int, error) {
	w.HideResponse = true
	w.Body = append(w.Body, bytes...)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.writeBody(string(response))
}

func (w *Writer) Clear() {
	w.Header().Del("Content-Type")
	w.Header().Del("Content-Length")
	w.Body = []byte{}
	w.Status = 0
}

func (w *Writer) Commit() (int, error) {
	w.Header().Set(
		"Content-Length",
		strconv.FormatInt(int64(len(w.Body)), 10),
	)
	if w.Status != 0 {
		w.ResponseWriter.WriteHeader(w.Status)
	}
	return w.ResponseWriter.Write(w.Body)
}

func (w *Writer) String() string {
	out := fmt.Sprintf("status %d (took %dms)", w.Status, w.Time)
	if len(w.Body) > 0 && !w.HideResponse {
		out = fmt.Sprintf("%s\n\tresponse: %s", out, w.Body)
	}
	return out
}
