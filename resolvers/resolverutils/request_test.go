package resolverutils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
		_, req, _ := CommonSetup(body)
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
	setup := func(user *database.User, params map[string]string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = SetContext(t, req, user, params)
		return req
	}

	t.Run("Normal", func(t *testing.T) {
		xUser := mocks.MakeUser()
		key1 := USERNAME_KEY
		key2 := CHAT_NAME_KEY
		contextParams := map[string]string{key1: "value1", key2: "value2"}
		req := setup(xUser, contextParams)

		user, routeParams, err := GetRequestContext(req, key1, key2)
		assert.IsNil(t, err)
		assert.DeepEquals(t, user, xUser)
		assert.DeepEquals(t, routeParams, &RouteParams{"value1", 0, "value2"})
	})

	t.Run("ParamExtractionFails", func(t *testing.T) {
		req := setup(mocks.MakeUser(), nil)

		_, _, err := GetRequestContext(req, USERNAME_KEY)
		assert.IsNotNil(t, err)
	})

	t.Run("NoUser", func(t *testing.T) {
		req := setup(nil, nil)
		buf := logger.MockFileLogger(t)

		_, _, err := GetRequestContext(req)
		assert.IsNotNil(t, err)
		assert.Equals(t, err.Status, http.StatusInternalServerError)
		assert.Equals(t, err.Message, "failed to fetch request user")
		assert.Contains(t, buf.String(), "[ERROR]", "user not found in request context")
	})
}
