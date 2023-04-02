package context

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
)

func TestGetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := mocks.MakeUser()
	req = req.WithContext(context.WithValue(req.Context(), userKey{}, xUser))

	user, _ := GetUser(req)
	assert.DeepEquals(t, user, xUser)
}

func TestGetUserButNotThere(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := GetUser(req)
	assert.ErrorHasMessage(t, err, "context user not found in request")
}

func TestGetUserButCastFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := struct{ Id string }{"asada"}
	req = req.WithContext(context.WithValue(req.Context(), userKey{}, xUser))

	_, err := GetUser(req)
	assert.ErrorHasMessage(t, err, "context user not of type User")
}

func TestSetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := mocks.MakeUser()
	req = SetUser(req, xUser)

	user := req.Context().Value(userKey{}).(*database.User)
	assert.DeepEquals(t, user, xUser)
}

func TestGetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	xValue := "testvalue"
	req = req.WithContext(context.WithValue(req.Context(), paramKey(key), xValue))

	value, _ := GetParam(req, key)
	assert.DeepEquals(t, value, xValue)
}

func TestGetParamButNotThere(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	req = req.WithContext(context.WithValue(req.Context(), paramKey("differekey"), "testvalue"))

	_, err := GetParam(req, key)
	assert.ErrorHasMessage(t, err, fmt.Sprintf("context parameter %s not found in request", key))
}

func TestGetParamButCastFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	value := struct{}{}
	req = req.WithContext(context.WithValue(req.Context(), paramKey(key), value))

	_, err := GetParam(req, key)
	assert.ErrorHasMessage(t, err, fmt.Sprintf("context parameter %s not of type string", key))
}

func TestSetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "foo"
	xValue := "bar"
	req = SetParam(req, key, xValue)

	value := req.Context().Value(paramKey(key)).(string)
	assert.DeepEquals(t, value, xValue)
}
