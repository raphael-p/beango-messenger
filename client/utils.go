package client

import (
	"fmt"
	"net/http"

	"github.com/raphael-p/beango/utils/response"
)

// TODO: unit test
func DisplayError(w *response.Writer, message string) {
	if message == "" {
		return
	}
	htmlStr := fmt.Sprintf("<div id='errors' class='error' hx-swap-oob='innerHTML'>%s</div>", message)
	w.WriteString(http.StatusOK, htmlStr)
}
