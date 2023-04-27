package context

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/raphael-p/beango/database"
)

// context keys, used to avoid clashes
type paramKey string
type userKey struct{}

func GetUser(r *http.Request) (*database.User, error) {
	rawUser := r.Context().Value(userKey{})
	if rawUser == nil {
		return nil, errors.New("user not found in request context")
	}
	user, ok := rawUser.(*database.User)
	if !ok {
		return nil, errors.New("user in request context not of type User")
	}
	if user == nil {
		return nil, errors.New("user in request context is nil")
	}
	return user, nil
}

func SetUser(r *http.Request, user *database.User) (*http.Request, error) {
	if user == nil {
		return r, errors.New("cannot set nil user to request context")
	}
	_, err := GetUser(r)
	if err == nil {
		return r, errors.New("user already in request context")
	}
	ctx := context.WithValue(r.Context(), userKey{}, user)
	return r.WithContext(ctx), nil
}

func GetParam(r *http.Request, key string) (string, error) {
	value := r.Context().Value(paramKey(key))
	if value == nil {
		return "", fmt.Errorf("path parameter %s not found", key)
	}
	stringValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("path parameter %s not of type string", key)
	}
	return stringValue, nil
}

func SetParam(r *http.Request, key string, value string) (*http.Request, error) {
	_, err := GetParam(r, key)
	if err == nil {
		return r, fmt.Errorf("path parameter %s already set", key)
	}
	ctx := context.WithValue(r.Context(), paramKey(key), value)
	return r.WithContext(ctx), nil
}
