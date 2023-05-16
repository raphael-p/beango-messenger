package resolvers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raphael-p/beango/test/assert"
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
		assert.Equals(t, w.Body, "expected `ptr` to be a pointer to a struct, got resolvers.TestStruct")
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
