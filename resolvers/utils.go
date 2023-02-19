package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raphael-p/beango/utils"
)

// Decodes JSON from HTTP request body and binds it to a struct pointer
func bindRequestJSON(w *utils.ResponseWriter, r *http.Request, ptr any) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(ptr); err != nil {
		fmt.Println(err.Error())
		w.StringResponse(http.StatusBadRequest, "malformed request body")
		return false
	}
	return true
}
