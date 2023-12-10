package response

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	if w.Status != 0 {
		w.ResponseWriter.WriteHeader(w.Status)
	}
}

func (w *Writer) writeBody(body []byte) (int, error) {
	w.Body = append(w.Body, body...)
	return w.ResponseWriter.Write(body)
}

func (w *Writer) WriteString(code int, response string) {
	w.WriteHeader(code)
	w.writeBody([]byte(response))
}

func (w *Writer) Write(bytes []byte) (int, error) {
	w.HideResponse = true
	if w.Status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.writeBody(bytes)
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
	w.writeBody(response)
}

func (w *Writer) WriteHTML(code int, response string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteString(code, response)
}

func (w *Writer) String() string {
	out := fmt.Sprintf("status %d (took %dms)", w.Status, w.Time)
	if len(w.Body) > 0 && !w.HideResponse {
		out = fmt.Sprintf("%s\n\tresponse: %s", out, w.Body)
	}
	return out
}
