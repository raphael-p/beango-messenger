package httputils

import (
	"fmt"
	"net/http"

	"github.com/raphael-p/beango/utils"
)

// context keys, used to avoid clashes
type ContextParameter string
type ContextUser string

func GetUserFromContext(r *http.Request) (any, error) {
	user := r.Context().Value(ContextUser("user"))
	if user == nil {
		message := "context user not found in request"
		utils.Logger.Error(message)
		return nil, fmt.Errorf(message)
	}
	return user, nil
}

func GetParamFromContext(r *http.Request, key string) (string, error) {
	value := r.Context().Value(ContextParameter(key))
	if value == nil {
		message := fmt.Sprintf("context parameter %s not found in request", key)
		utils.Logger.Error(message)
		return "", fmt.Errorf(message)
	}
	stringValue, ok := value.(string)
	if !ok {
		message := fmt.Sprintf("context parameter %s not of type string", key)
		utils.Logger.Error(message)
		return "", fmt.Errorf(message)
	}
	return stringValue, nil
}
