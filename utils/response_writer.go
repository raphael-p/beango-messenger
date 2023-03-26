package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseWriter struct {
	Status int
	Body   string
	Time   int64
	http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w}
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriter) Write(body []byte) (int, error) {
	w.Body = string(body)
	return w.ResponseWriter.Write(body)
}

func (w *ResponseWriter) StringResponse(code int, response string) {
	w.WriteHeader(code)
	w.Write([]byte(response))
}

func (w *ResponseWriter) JSONResponse(code int, responseObject any) {
	w.WriteHeader(code)
	response, err := json.Marshal(responseObject)
	if err != nil {
		w.StringResponse(http.StatusBadRequest, err.Error())
	}
	w.Header().Set("content-type", "application/json")
	w.Write(response)
}

func (w *ResponseWriter) String() string {
	out := fmt.Sprintf("status %d (took %dms)", w.Status, w.Time)
	if w.Body != "" {
		out = fmt.Sprintf("%s\n\tresponse: %s", out, w.Body)
	}
	return out
}
