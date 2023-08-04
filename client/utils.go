package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/raphael-p/beango/utils/response"
)

// TODO: unit test
func displayError(w *response.Writer, message string) {
	if message == "" {
		return
	}
	htmlStr := fmt.Sprintf("<div id='errors' class='error' hx-swap-oob='innerHTML'>%s</div>", message)
	w.WriteString(http.StatusOK, htmlStr)
}

// TODO: unit test
func cloneRequest(r *http.Request) *http.Request {
	clone := r.Clone(r.Context())
	// TODO: investigate error handling: can ignoring this error cause a panic?
	if body, err := io.ReadAll(r.Body); err == nil {
		clone.Body = io.NopCloser(bytes.NewReader(body))
	}
	return clone
}
