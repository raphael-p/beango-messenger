package resolverutils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/utils/logger"
)

func TestExtractRouteParams(t *testing.T) {
	setup := func(params map[string]string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = SetContext(t, req, nil, params)
		return req
	}

	t.Run("GetAllParams", func(t *testing.T) {
		key1 := USERNAME_KEY
		key2 := CHAT_NAME_KEY
		key3 := CHAT_ID_KEY
		contextParams := map[string]string{key1: "value1", key2: "value2", key3: "29"}
		req := setup(contextParams)

		routeParams, err := extractRouteParams(req, key1, key2)
		assert.IsNil(t, err)
		assert.DeepEquals(t, routeParams, &RouteParams{"value1", 0, "value2"})
	})

	t.Run("ExtraRequestParam", func(t *testing.T) {
		key := USERNAME_KEY
		extraKey := CHAT_NAME_KEY
		contextParams := map[string]string{key: "value1", extraKey: "value2"}
		req := setup(contextParams)

		params, err := extractRouteParams(req, key)
		assert.IsNil(t, err)
		assert.DeepEquals(t, params, &RouteParams{"value1", 0, ""})
	})

	t.Run("NoParamKeysPassed", func(t *testing.T) {
		req := setup(map[string]string{"param1": "value1"})

		params, err := extractRouteParams(req)
		assert.IsNil(t, err)
		assert.DeepEquals(t, params, &RouteParams{})
	})

	t.Run("MissingParamInRequest", func(t *testing.T) {
		key1 := USERNAME_KEY
		key2 := CHAT_NAME_KEY
		req := setup(map[string]string{key1: "some-value"})
		buf := logger.MockFileLogger(t)

		_, err := extractRouteParams(req, key1, key2)
		assert.IsNotNil(t, err)
		assert.Equals(t, err.Status, http.StatusInternalServerError)
		assert.Equals(t, err.Message, fmt.Sprint("failed to fetch path parameter: ", key2))
		assert.Contains(t, buf.String(), "[ERROR]", fmt.Sprintf("path parameter %s not found", key2))
	})

	t.Run("ChatIDNotAnInt", func(t *testing.T) {
		key := CHAT_ID_KEY
		req := setup(map[string]string{key: "some-value"})

		_, err := extractRouteParams(req, key)
		assert.IsNotNil(t, err)
		assert.Equals(t, err.Status, http.StatusBadRequest)
		assert.Equals(t, err.Message, "chat ID must be an integer")
	})
}
