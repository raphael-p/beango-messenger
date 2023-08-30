package resolverutils

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/logger"
)

func TestBindRequestJSON(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	setup := func(body string, ptr any) *HTTPError {
		_, req := MockRequest(body)
		httpError := GetRequestBody(req, ptr)
		return httpError
	}

	t.Run("Normal", func(t *testing.T) {
		name := "John"
		age := 30
		body := fmt.Sprintf(`{"name": "%s", "Age": %d}`, name, age)
		var testStruct TestStruct
		err := setup(body, &testStruct)
		assert.IsNil(t, err)
		xTestStruct := TestStruct{name, age}
		assert.DeepEquals(t, testStruct, xTestStruct)
	})

	t.Run("NonPointerBind", func(t *testing.T) {
		var testStruct TestStruct
		err := setup(`{"name": "John", "Age": 30}`, testStruct)
		assert.IsNotNil(t, err)
		assert.Equals(t, err.Status, http.StatusBadRequest)
		xMessage := "expected `ptr` to be a pointer to a struct, got resolverutils.TestStruct"
		assert.Equals(t, err.Message, xMessage)
	})

	t.Run("NonStructBind", func(t *testing.T) {
		var testStruct *string
		err := setup(`{"name": "John", "Age": 30}`, testStruct)
		assert.IsNotNil(t, err)
		assert.Equals(t, err.Status, http.StatusBadRequest)
		assert.Equals(t, err.Message, "expected `ptr` to be a pointer to a struct, got *string")
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		var testStruct TestStruct
		err := setup(`{"name": "John"}`, &testStruct)
		assert.IsNotNil(t, err)
		assert.Equals(t, err.Status, http.StatusBadRequest)
		assert.Equals(t, err.Message, "missing required field(s): [age]")
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		var testStruct TestStruct
		err := setup(`{"name": "John",}`, &testStruct)
		assert.IsNotNil(t, err)
		assert.Equals(t, err.Status, http.StatusBadRequest)
		assert.Contains(t, err.Message, "malformed request body: ")
	})
}

func TestGetRequestContext(t *testing.T) {
	setup := func(user *database.User, params map[string]string) (*http.Request, *bytes.Buffer) {
		_, req := MockRequest("")
		req = SetContext(t, req, user, params)
		buf := logger.MockFileLogger(t)
		return req, buf
	}

	t.Run("Normal", func(t *testing.T) {
		xUser := mocks.MakeUser()
		key1 := USERNAME_KEY
		key2 := CHAT_NAME_KEY
		contextParams := map[string]string{key1: "value1", key2: "value2"}
		req, _ := setup(xUser, contextParams)

		user, routeParams, err := GetRequestContext(req, key1, key2)
		assert.IsNil(t, err)
		assert.DeepEquals(t, user, xUser)
		assert.DeepEquals(t, routeParams, &RouteParams{"value1", 0, "value2"})
	})

	t.Run("ExtraRequestParam", func(t *testing.T) {
		key := USERNAME_KEY
		extraKey := CHAT_NAME_KEY
		contextParams := map[string]string{key: "value1", extraKey: "value2"}
		req, _ := setup(mocks.MakeUser(), contextParams)

		_, params, err := GetRequestContext(req, key)
		assert.IsNil(t, err)
		assert.DeepEquals(t, params, &RouteParams{"value1", 0, ""})
	})

	t.Run("NoParamKeys", func(t *testing.T) {
		req, _ := setup(mocks.MakeUser(), map[string]string{"param1": "value1"})

		_, params, err := GetRequestContext(req)
		assert.IsNil(t, err)
		assert.DeepEquals(t, params, &RouteParams{})
	})

	t.Run("MissingRequestParam", func(t *testing.T) {
		key1 := USERNAME_KEY
		key2 := CHAT_NAME_KEY
		req, buf := setup(mocks.MakeUser(), map[string]string{key1: "some-value"})

		_, _, err := GetRequestContext(req, key1, key2)
		assert.IsNotNil(t, err)
		assert.Equals(t, err.Status, http.StatusInternalServerError)
		assert.Equals(t, err.Message, fmt.Sprint("failed to fetch path parameter: ", key2))
		assert.Contains(t, buf.String(), "[ERROR]", fmt.Sprintf("path parameter %s not found", key2))
	})

	t.Run("NoUser", func(t *testing.T) {
		req, buf := setup(nil, nil)

		_, _, err := GetRequestContext(req)
		assert.IsNotNil(t, err)
		assert.Equals(t, err.Status, http.StatusInternalServerError)
		assert.Equals(t, err.Message, "failed to fetch request user")
		assert.Contains(t, buf.String(), "[ERROR]", "user not found in request context")
	})
}
