package utils

import (
	"net/http"

	"github.com/raphael-p/beango/database"
)

// context keys, used to avoid clashes
type ContextParameters string
type ContextUser string

func GetUserFromContext(r *http.Request) *database.User {
	return r.Context().Value(ContextUser("user")).(*database.User)
}

func GetParamFromContext(r *http.Request, key string) string {
	return r.Context().Value(ContextParameters(key)).(string)
}
