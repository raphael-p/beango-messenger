package httputils_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/raphael-p/beango/httputils"
)

func TestGetUser(t *testing.T) {
	xUser := struct {
		ID       string
		Username string
	}{
		ID:       "123",
		Username: "johndoe",
	}
	req := &http.Request{Header: http.Header{}}
	req = req.WithContext(context.WithValue(req.Context(), httputils.ContextUser("user"), xUser))

	if user, _ := httputils.GetUserFromContext(req); !reflect.DeepEqual(xUser, user) {
		t.Errorf("expected user %v, but got %v", xUser, user)
	}
}

func TestGetUserButNotThere(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
	}
	xError := fmt.Errorf("context user not found in request")

	if _, err := httputils.GetUserFromContext(req); xError == err {
		t.Errorf("expected error %v, but got %v", xError, err)
	}
}

func TestGetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	key := "testkey"
	xValue := "testvalue"
	req = req.WithContext(context.WithValue(req.Context(), httputils.ContextParameter(key), xValue))

	if value, _ := httputils.GetParamFromContext(req, key); xValue != value {
		t.Errorf("expected %v, but got %v", xValue, value)

	}
}

func TestGetParamButNotThere(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
	}
	key := "testkey"
	xError := fmt.Errorf("context parameter %s not found in request", key)

	if _, err := httputils.GetParamFromContext(req, key); xError == err {
		t.Errorf("expected error %v, but got %v", xError, err)
	}
}

func TestGetParamButCastFails(t *testing.T) {
	req := &http.Request{
		Header: http.Header{},
	}
	key := "testkey"
	value := struct{}{}
	req = req.WithContext(context.WithValue(req.Context(), httputils.ContextParameter(key), value))
	xError := fmt.Errorf("context parameter %s not of type string", key)

	if _, err := httputils.GetParamFromContext(req, key); xError == err {
		t.Errorf("expected error %v, but got %v", xError, err)
	}
}
