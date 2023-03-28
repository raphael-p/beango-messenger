package httputils

import (
	"net/http"
)

// context keys, used to avoid clashes
type ContextParameter string
type ContextUser string

func GetUserFromContext(r *http.Request) any {
	return r.Context().Value(ContextUser("user"))
}

func GetParamFromContext(r *http.Request, key string) string {
	return r.Context().Value(ContextParameter(key)).(string)
}
