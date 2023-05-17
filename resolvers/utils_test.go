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
	type param = struct{ key, value string }
	setup := func(user *database.User, params []param) (
		*http.Request,
		*response.Writer,
		*bytes.Buffer,
	) {
		req := httptest.NewRequest("GET", "/test", nil)
		if user != nil {
			newReq, err := context.SetUser(req, user)
			assert.IsNil(t, err)
			req = newReq
		}
		for _, param := range params {
			newReq, err := context.SetParam(req, param.key, param.value)
			assert.IsNil(t, err)
			req = newReq
		}
		w := response.NewWriter(httptest.NewRecorder())
		buf := logger.MockFileLogger(t)
		return req, w, buf
	}

	t.Run("Normal", func(t *testing.T) {
		xUser := mocks.MakeUser()
		param1 := "value1"
		param2 := "value2"
		params := []param{{"Param1", param1}, {"Param2", param2}}
		req, _, _ := setup(xUser, params)

		type TestStruct struct{ Param1, Param2 string }
		var testStruct TestStruct
		user, ok := getRequestContext(nil, req, &testStruct)
		assert.Equals(t, ok, true)
		assert.DeepEquals(t, user, xUser)
		assert.DeepEquals(t, testStruct, TestStruct{param1, param2})
	})

	t.Run("ParamMissingInPointer", func(t *testing.T) {
		param1 := "value1"
		params := []param{{"Param1", param1}, {"Param2", "value2"}}
		req, _, _ := setup(mocks.MakeUser(), params)
		type TestStruct struct{ Param1 string }
		var testStruct TestStruct

		_, ok := getRequestContext(nil, req, &testStruct)
		assert.Equals(t, ok, true)
		assert.DeepEquals(t, testStruct, TestStruct{param1})
	})

	t.Run("NoUser", func(t *testing.T) {
		req, w, buf := setup(nil, nil)

		_, ok := getRequestContext(w, req, nil)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, w.Body, "failed to fetch user")
		assert.Contains(t, buf.String(), "[ERROR]", "user not found in request context")
	})

	t.Run("InvalidPointer", func(t *testing.T) {
		req, w, buf := setup(mocks.MakeUser(), nil)
		type TestStruct struct {
			Param1 string
			Param2 int
		}

		_, ok := getRequestContext(w, req, &TestStruct{})
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, w.Body, "failed to fetch path parameters")
		xMessage := "path param variable must point to a struct of strings"
		assert.Contains(t, buf.String(), "[ERROR]", xMessage)
	})

	t.Run("ParamMissingInRequest", func(t *testing.T) {
		req, w, buf := setup(mocks.MakeUser(), []param{{"Param1", "value1"}})
		type TestStruct struct{ Param1, Param2 string }

		_, ok := getRequestContext(w, req, &TestStruct{})
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, w.Body, "failed to fetch path parameter: Param2")
		assert.Contains(t, buf.String(), "[ERROR]", "path parameter Param2 not found")
	})
}
