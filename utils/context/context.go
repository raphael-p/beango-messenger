package context

import (
	"context"
	"fmt"
	"net/http"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/logger"
)

// context keys, used to avoid clashes
type paramKey string
type userKey struct{}

func GetUser(r *http.Request) (*database.User, error) {
	rawUser := r.Context().Value(userKey{})
	if rawUser == nil {
		message := "context user not found in request"
		logger.Error(message)
		return nil, fmt.Errorf(message)
	}
	user, ok := rawUser.(*database.User)
	if !ok {
		message := "context user not of type User"
		logger.Error(message)
		return nil, fmt.Errorf(message)
	}
	return user, nil
}

func SetUser(r *http.Request, user *database.User) *http.Request {
	ctx := context.WithValue(r.Context(), userKey{}, user)
	return r.WithContext(ctx)
}

func GetParam(r *http.Request, key string) (string, error) {
	value := r.Context().Value(paramKey(key))
	if value == nil {
		message := fmt.Sprintf("context parameter %s not found in request", key)
		logger.Error(message)
		return "", fmt.Errorf(message)
	}
	stringValue, ok := value.(string)
	if !ok {
		message := fmt.Sprintf("context parameter %s not of type string", key)
		logger.Error(message)
		return "", fmt.Errorf(message)
	}
	return stringValue, nil
}

func SetParam(r *http.Request, key string, value string) *http.Request {
	ctx := context.WithValue(r.Context(), paramKey(key), value)
	return r.WithContext(ctx)
}
