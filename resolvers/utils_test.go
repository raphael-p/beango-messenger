package resolvers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func TestBindRequestJSON(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	setup := func(body string, ptr any) (bool, *response.Writer) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := response.NewWriter(httptest.NewRecorder())
		ok := bindRequestJSON(w, req, ptr)
		return ok, w
	}

	t.Run("Normal", func(t *testing.T) {
		name := "John"
		age := 30
		body := fmt.Sprintf(`{"name": "%s", "Age": %d}`, name, age)
		var testStruct TestStruct
		ok, _ := setup(body, &testStruct)
		assert.Equals(t, ok, true)
		xTestStruct := TestStruct{name, age}
		assert.DeepEquals(t, testStruct, xTestStruct)
	})

	t.Run("NonPointerBind", func(t *testing.T) {
		var testStruct TestStruct
		ok, w := setup(`{"name": "John", "Age": 30}`, testStruct)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		xBody := "expected `ptr` to be a pointer to a struct, got resolvers.TestStruct"
		assert.Equals(t, w.Body, xBody)
	})

	t.Run("NonStructBind", func(t *testing.T) {
		var testStruct *string
		ok, w := setup(`{"name": "John", "Age": 30}`, testStruct)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Equals(t, w.Body, "expected `ptr` to be a pointer to a struct, got *string")
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		var testStruct TestStruct
		ok, w := setup(`{"name": "John"}`, &testStruct)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Equals(t, w.Body, "missing required field(s): [age]")
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		var testStruct TestStruct
		ok, w := setup(`{"name": "John",}`, &testStruct)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Contains(t, w.Body, "malformed request body: ")
	})
}

func TestGetRequestContext(t *testing.T) {
	setup := func(user *database.User, params map[string]string) (
		*http.Request,
		*response.Writer,
		*bytes.Buffer,
	) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		if user != nil {
			newReq, err := context.SetUser(req, user)
			assert.IsNil(t, err)
			req = newReq
		}
		for key, value := range params {
			newReq, err := context.SetParam(req, key, value)
			assert.IsNil(t, err)
			req = newReq
		}
		w := response.NewWriter(httptest.NewRecorder())
		buf := logger.MockFileLogger(t)
		return req, w, buf
	}

	t.Run("Normal", func(t *testing.T) {
		xUser := mocks.MakeUser()
		key1 := "Param1"
		key2 := "param2"
		xParams := map[string]string{key1: "value1", key2: "value2"}
		req, _, _ := setup(xUser, xParams)

		user, params, ok := getRequestContext(nil, req, []string{key1, key2})
		assert.Equals(t, ok, true)
		assert.DeepEquals(t, user, xUser)
		assert.DeepEquals(t, params, xParams)
	})

	t.Run("ExtraRequestParam", func(t *testing.T) {
		key := "param1"
		extraKey := "extra"
		xParams := map[string]string{key: "value1", extraKey: "value2"}
		req, _, _ := setup(mocks.MakeUser(), xParams)

		_, params, ok := getRequestContext(nil, req, []string{key})
		assert.Equals(t, ok, true)
		delete(xParams, extraKey)
		assert.DeepEquals(t, params, xParams)
	})

	t.Run("NilKeySlice", func(t *testing.T) {
		req, _, _ := setup(mocks.MakeUser(), map[string]string{"param1": "value1"})

		_, params, ok := getRequestContext(nil, req, nil)
		assert.Equals(t, ok, true)
		assert.DeepEquals(t, params, map[string]string{})
	})

	t.Run("MissingRequestParam", func(t *testing.T) {
		key1 := "param1"
		key2 := "param2"
		req, w, buf := setup(mocks.MakeUser(), map[string]string{key1: "some-value"})

		_, _, ok := getRequestContext(w, req, []string{key1, key2})
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, w.Body, fmt.Sprint("failed to fetch path parameter: ", key2))
		assert.Contains(t, buf.String(), "[ERROR]", fmt.Sprintf("path parameter %s not found", key2))
	})

	t.Run("NoUser", func(t *testing.T) {
		req, w, buf := setup(nil, nil)

		_, _, ok := getRequestContext(w, req, []string{})
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, w.Body, "failed to fetch user")
		assert.Contains(t, buf.String(), "[ERROR]", "user not found in request context")
	})
}
