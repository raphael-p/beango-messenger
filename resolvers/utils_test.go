package resolvers

import (
	"bytes"
	"errors"
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

func mockRequest(body string) (*response.Writer, *http.Request) {
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := response.NewWriter(httptest.NewRecorder())
	return w, req
}

func setContext(
	t *testing.T,
	req *http.Request,
	user *database.User,
	params map[string]string,
) *http.Request {
	var err error = nil
	if user != nil {
		req, err = context.SetUser(req, user)
		assert.IsNil(t, err)
	}
	for key, value := range params {
		req, err = context.SetParam(req, key, value)
		assert.IsNil(t, err)
	}
	return req
}

func TestBindRequestJSON(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	setup := func(body string, ptr any) (bool, *response.Writer) {
		w, req := mockRequest(body)
		ok := getRequestBody(w, req, ptr)
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
		assert.Equals(t, string(w.Body), xBody)
	})

	t.Run("NonStructBind", func(t *testing.T) {
		var testStruct *string
		ok, w := setup(`{"name": "John", "Age": 30}`, testStruct)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Equals(t, string(w.Body), "expected `ptr` to be a pointer to a struct, got *string")
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		var testStruct TestStruct
		ok, w := setup(`{"name": "John"}`, &testStruct)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Equals(t, string(w.Body), "missing required field(s): [age]")
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		var testStruct TestStruct
		ok, w := setup(`{"name": "John",}`, &testStruct)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Contains(t, string(w.Body), "malformed request body: ")
	})
}

func TestGetRequestContext(t *testing.T) {
	setup := func(user *database.User, params map[string]string) (
		*http.Request,
		*response.Writer,
		*bytes.Buffer,
	) {
		w, req := mockRequest("")
		req = setContext(t, req, user, params)
		buf := logger.MockFileLogger(t)
		return req, w, buf
	}

	t.Run("Normal", func(t *testing.T) {
		xUser := mocks.MakeUser()
		key1 := "Param1"
		key2 := "param2"
		xParams := map[string]string{key1: "value1", key2: "value2"}
		req, _, _ := setup(xUser, xParams)

		user, params, ok := getRequestContext(nil, req, key1, key2)
		assert.Equals(t, ok, true)
		assert.DeepEquals(t, user, xUser)
		assert.DeepEquals(t, params, xParams)
	})

	t.Run("ExtraRequestParam", func(t *testing.T) {
		key := "param1"
		extraKey := "extra"
		xParams := map[string]string{key: "value1", extraKey: "value2"}
		req, _, _ := setup(mocks.MakeUser(), xParams)

		_, params, ok := getRequestContext(nil, req, key)
		assert.Equals(t, ok, true)
		delete(xParams, extraKey)
		assert.DeepEquals(t, params, xParams)
	})

	t.Run("NoParamKeys", func(t *testing.T) {
		req, _, _ := setup(mocks.MakeUser(), map[string]string{"param1": "value1"})

		_, params, ok := getRequestContext(nil, req)
		assert.Equals(t, ok, true)
		assert.DeepEquals(t, params, map[string]string{})
	})

	t.Run("MissingRequestParam", func(t *testing.T) {
		key1 := "param1"
		key2 := "param2"
		req, w, buf := setup(mocks.MakeUser(), map[string]string{key1: "some-value"})

		_, _, ok := getRequestContext(w, req, key1, key2)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, string(w.Body), fmt.Sprint("failed to fetch path parameter: ", key2))
		assert.Contains(t, buf.String(), "[ERROR]", fmt.Sprintf("path parameter %s not found", key2))
	})

	t.Run("NoUser", func(t *testing.T) {
		req, w, buf := setup(nil, nil)

		_, _, ok := getRequestContext(w, req)
		assert.Equals(t, ok, false)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, string(w.Body), "failed to fetch request user")
		assert.Contains(t, buf.String(), "[ERROR]", "user not found in request context")
	})
}

func TestHandleDatabaseError(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		buf := logger.MockFileLogger(t)
		errPrefix := "database operation failed"
		errMessage := "this did not go well"
		httpError := HandleDatabaseError(errors.New(errMessage))
		assert.Equals(t, httpError.status, http.StatusInternalServerError)
		assert.Equals(t, httpError.message, errPrefix)
		assert.Contains(t, buf.String(), "[ERROR] "+errPrefix+": "+errMessage)
	})
}
