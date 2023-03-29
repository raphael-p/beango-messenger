package httputils_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/raphael-p/beango/httputils"
	"github.com/raphael-p/beango/test/utils/database"
)

func TestGetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := database.MakeUser()
	req = req.WithContext(context.WithValue(req.Context(), httputils.ContextUser("user"), xUser))
	if user, _ := httputils.GetContextUser(req); !reflect.DeepEqual(xUser, user) {
		t.Errorf("expected user %v, but got %v", xUser, user)
	}
}

func TestGetUserButNotThere(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xMessage := "context user not found in request"

	if _, err := httputils.GetContextUser(req); err.Error() != xMessage {
		t.Errorf("expected error %v, but got %v", xMessage, err)
	}
}

func TestGetUserButCastFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	xUser := struct{ Id string }{"asada"}
	req = req.WithContext(context.WithValue(req.Context(), httputils.ContextUser("user"), xUser))
	xMessage := "context user not of type User"

	if _, err := httputils.GetContextUser(req); err.Error() != xMessage {
		t.Errorf("expected error %v, but got %v", xMessage, err)
	}
}

func TestGetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	xValue := "testvalue"
	req = req.WithContext(context.WithValue(req.Context(), httputils.ContextParameter(key), xValue))

	if value, _ := httputils.GetContextParam(req, key); value != xValue {
		t.Errorf("expected %v, but got %v", xValue, value)

	}
}

func TestGetParamButNotThere(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	xMessage := fmt.Sprintf("context parameter %s not found in request", key)

	if _, err := httputils.GetContextParam(req, key); err.Error() != xMessage {
		t.Errorf("expected error %v, but got %v", xMessage, err)
	}
}

func TestGetParamButCastFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	value := struct{}{}
	req = req.WithContext(context.WithValue(req.Context(), httputils.ContextParameter(key), value))
	xMessage := fmt.Sprintf("context parameter %s not of type string", key)

	if _, err := httputils.GetContextParam(req, key); err.Error() != xMessage {
		t.Errorf("expected error %v, but got %v", xMessage, err)
	}
}
