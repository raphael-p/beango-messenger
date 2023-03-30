package httputils

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raphael-p/beango/database"
	assert "github.com/raphael-p/beango/test/assertions"
	"github.com/raphael-p/beango/test/mocks"
)

func TestGetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := mocks.MakeUser()
	req = req.WithContext(context.WithValue(req.Context(), ContextUser("user"), xUser))

	user, _ := GetContextUser(req)
	assert.DeepEquals(t, user, xUser)
}

func TestGetUserButNotThere(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := GetContextUser(req)
	assert.ErrorHasMessage(t, err, "context user not found in request")
}

func TestGetUserButCastFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := struct{ Id string }{"asada"}
	req = req.WithContext(context.WithValue(req.Context(), ContextUser("user"), xUser))

	_, err := GetContextUser(req)
	assert.ErrorHasMessage(t, err, "context user not of type User")
}

func TestSetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := mocks.MakeUser()
	req = SetContextUser(req, xUser)

	user := req.Context().Value(ContextUser("user")).(*database.User)
	assert.DeepEquals(t, user, xUser)
}

func TestGetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	xValue := "testvalue"
	req = req.WithContext(context.WithValue(req.Context(), ContextParameter(key), xValue))

	value, _ := GetContextParam(req, key)
	assert.DeepEquals(t, value, xValue)
}

func TestGetParamButNotThere(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"

	_, err := GetContextParam(req, key)
	assert.ErrorHasMessage(t, err, fmt.Sprintf("context parameter %s not found in request", key))
}

func TestGetParamButCastFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	value := struct{}{}
	req = req.WithContext(context.WithValue(req.Context(), ContextParameter(key), value))

	_, err := GetContextParam(req, key)
	assert.ErrorHasMessage(t, err, fmt.Sprintf("context parameter %s not of type string", key))
}

func TestSetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "foo"
	xValue := "bar"
	req = SetContextParam(req, key, xValue)

	value := req.Context().Value(ContextParameter(key)).(string)
	assert.DeepEquals(t, value, xValue)
}
