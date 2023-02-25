package resolvers

import (
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w *utils.ResponseWriter, r *http.Request) {
	var input LoginInput
	if ok := bindRequestJSON(w, r, &input); !ok {
		return
	}

	user, ok := database.GetUserByUsername(input.Username)
	if ok {
		err := bcrypt.CompareHashAndPassword(user.Key, []byte(input.Password))
		if err == nil {
			w.StringResponse(http.StatusOK, "token")
			return
		}
	}

	w.StringResponse(http.StatusUnauthorized, "authentication failed")
}
