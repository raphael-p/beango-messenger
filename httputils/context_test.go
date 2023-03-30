package httputils

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/utils/mocks"
)

func TestGetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := mocks.MakeUser()
	req = req.WithContext(context.WithValue(req.Context(), ContextUser("user"), xUser))

	if user, _ := GetContextUser(req); !reflect.DeepEqual(xUser, user) {
		t.Errorf("expected user %v, but got %v", xUser, user)
	}
}

func TestGetUserButNotThere(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	xMessage := "context user not found in request"
	if _, err := GetContextUser(req); err.Error() != xMessage {
		t.Errorf("expected error %v, but got %v", xMessage, err)
	}
}

func TestGetUserButCastFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := struct{ Id string }{"asada"}
	req = req.WithContext(context.WithValue(req.Context(), ContextUser("user"), xUser))

	xMessage := "context user not of type User"
	if _, err := GetContextUser(req); err.Error() != xMessage {
		t.Errorf("expected error %v, but got %v", xMessage, err)
	}
}

func TestSetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := mocks.MakeUser()
	req = SetContextUser(req, xUser)

	user := req.Context().Value(ContextUser("user")).(*database.User)
	if !reflect.DeepEqual(user, xUser) {
		t.Errorf("expected user %v, but got %v", xUser, user)
	}
}

func TestGetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	xValue := "testvalue"
	req = req.WithContext(context.WithValue(req.Context(), ContextParameter(key), xValue))

	if value, _ := GetContextParam(req, key); value != xValue {
		t.Errorf("expected %v, but got %v", xValue, value)

	}
}

func TestGetParamButNotThere(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"

	xMessage := fmt.Sprintf("context parameter %s not found in request", key)
	if _, err := GetContextParam(req, key); err.Error() != xMessage {
		t.Errorf("expected error %v, but got %v", xMessage, err)
	}
}

func TestGetParamButCastFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	value := struct{}{}
	req = req.WithContext(context.WithValue(req.Context(), ContextParameter(key), value))

	xMessage := fmt.Sprintf("context parameter %s not of type string", key)
	if _, err := GetContextParam(req, key); err.Error() != xMessage {
		t.Errorf("expected error %v, but got %v", xMessage, err)
	}
}

func TestSetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "foo"
	xValue := "bar"
	req = SetContextParam(req, key, xValue)

	value := req.Context().Value(ContextParameter(key)).(string)
	if value != xValue {
		t.Errorf("expected %v, but got %v", xValue, value)
	}
}
