package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func displayError(w *response.Writer, message string) {
	if message == "" {
		return
	}
	htmlStr := fmt.Sprintf("<div id='errors' hx-swap-oob='true'>%s</div>", message)
	w.WriteString(http.StatusOK, htmlStr)
}

func cloneRequest(w *response.Writer, r *http.Request, count int) ([]*http.Request, bool) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("error reading request body: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}

	clones := make([]*http.Request, count)
	for i := 0; i < count; i++ {
		var rClone *http.Request
		if i == 0 {
			rClone = r
		} else {
			rClone = r.Clone(r.Context())
		}
		rClone.Body = io.NopCloser(bytes.NewReader(body))
		clones[i] = rClone
	}

	return clones, true
}
