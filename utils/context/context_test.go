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

	user, err := GetUser(req)
	assert.IsNil(t, err)
	assert.DeepEquals(t, user, xUser)
}

func TestGetMultipleUsers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := mocks.MakeUser()
	req = req.WithContext(context.WithValue(req.Context(), userKey{}, mocks.MakeUser()))
	req = req.WithContext(context.WithValue(req.Context(), userKey{}, xUser))

	user, err := GetUser(req)
	assert.IsNil(t, err)
	assert.DeepEquals(t, user, xUser)
}

func TestGetUserButNotThere(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := GetUser(req)
	assert.ErrorHasMessage(t, err, "context user not found in request")
}

func TestGetUserButCastFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := struct{ ID string }{"a-unique-id"}
	req = req.WithContext(context.WithValue(req.Context(), userKey{}, xUser))

	_, err := GetUser(req)
	assert.ErrorHasMessage(t, err, "context user not of type User")
}

func TestSetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := mocks.MakeUser()

	req, err := SetUser(req, xUser)
	assert.IsNil(t, err)
	user := req.Context().Value(userKey{}).(*database.User)
	assert.DeepEquals(t, user, xUser)
}

func TestSetMultipleUsers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := mocks.MakeUser()

	req, err := SetUser(req, xUser)
	assert.IsNil(t, err)
	req, err = SetUser(req, mocks.MakeUser())
	assert.ErrorHasMessage(t, err, "context user already set")
	user := req.Context().Value(userKey{}).(*database.User)
	assert.DeepEquals(t, user, xUser)
}

func TestGetParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key1 := "testkey1"
	key2 := "testkey2"
	xValue1 := "testvalue1"
	xValue2 := "testvalue2"
	req = req.WithContext(context.WithValue(req.Context(), paramKey(key1), xValue1))
	req = req.WithContext(context.WithValue(req.Context(), paramKey(key2), xValue2))

	value1, err := GetParam(req, key1)
	assert.IsNil(t, err)
	value2, err := GetParam(req, key2)
	assert.IsNil(t, err)
	assert.Equals(t, value1, xValue1)
	assert.Equals(t, value2, xValue2)
}

func TestGetSameParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	xValue := "testvalue2"
	req = req.WithContext(context.WithValue(req.Context(), paramKey(key), "testvalue1"))
	req = req.WithContext(context.WithValue(req.Context(), paramKey(key), xValue))

	value1, err := GetParam(req, key)
	assert.IsNil(t, err)
	value2, err := GetParam(req, key)
	assert.IsNil(t, err)
	assert.Equals(t, value1, xValue)
	assert.Equals(t, value2, xValue)
}

func TestGetParamCalledUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "user"
	xValue := "testvalue"
	req = req.WithContext(context.WithValue(req.Context(), paramKey(key), xValue))
	req = req.WithContext(context.WithValue(req.Context(), userKey{}, mocks.MakeUser()))

	value, err := GetParam(req, key)
	assert.IsNil(t, err)
	assert.Equals(t, value, xValue)
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
	key1 := "foo"
	key2 := "secondkey"
	xValue1 := "bar"
	xValue2 := "secondvalue"

	req, err := SetParam(req, key1, xValue1)
	assert.IsNil(t, err)
	req, err = SetParam(req, key2, xValue2)
	assert.IsNil(t, err)
	value1 := req.Context().Value(paramKey(key1)).(string)
	value2 := req.Context().Value(paramKey(key2)).(string)
	assert.Equals(t, value1, xValue1)
	assert.Equals(t, value2, xValue2)
}

func TestSetSameParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "samekey"
	xValue := "test-value-1"

	req, err := SetParam(req, key, xValue)
	assert.IsNil(t, err)
	req, err = SetParam(req, key, "test-value-2")
	assert.ErrorHasMessage(t, err, fmt.Sprintf("context parameter %s already set", key))
	value := req.Context().Value(paramKey(key)).(string)
	assert.Equals(t, value, xValue)
}
