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
		message := "context user not found in request"
		return nil, errors.New(message)
	}
	user, ok := rawUser.(*database.User)
	if !ok {
		message := "context user not of type User"
		return nil, errors.New(message)
	}
	return user, nil
}

func SetUser(r *http.Request, user *database.User) (*http.Request, error) {
	_, err := GetUser(r)
	if err == nil {
		return r, errors.New("context user already set")
	}
	ctx := context.WithValue(r.Context(), userKey{}, user)
	return r.WithContext(ctx), nil
}

func GetParam(r *http.Request, key string) (string, error) {
	value := r.Context().Value(paramKey(key))
	if value == nil {
		return "", fmt.Errorf("context parameter %s not found in request", key)

	}
	stringValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("context parameter %s not of type string", key)

	}
	return stringValue, nil
}

func SetParam(r *http.Request, key string, value string) (*http.Request, error) {
	_, err := GetParam(r, key)
	if err == nil {
		return r, fmt.Errorf("context parameter %s already set", key)
	}
	ctx := context.WithValue(r.Context(), paramKey(key), value)
	return r.WithContext(ctx), nil
}
