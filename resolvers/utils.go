package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raphael-p/beango/utils"
)

func bindJSON(w *utils.ResponseWriter, r *http.Request, input any) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(input); err != nil {
		fmt.Println(err.Error())
		w.StringResponse(http.StatusBadRequest, "malformed request body")
		return false
	}
	return true
}
