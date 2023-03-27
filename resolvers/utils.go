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
		response := fmt.Sprint("malformed request body: ", err)
		w.StringResponse(http.StatusBadRequest, response)
		return false
	}
	if fields := utils.ValidateRequiredFields(ptr); len(fields) != 0 {
		response := fmt.Sprintf("missing required field(s): %s", fields)
		w.StringResponse(http.StatusBadRequest, response)
		return false
	}
	return true
}
