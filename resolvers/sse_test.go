package resolvers

import (
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/utils/collections"
	"github.com/raphael-p/beango/utils/response"
)

func TestRegisterConnection(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		var key int64 = 1
		xStatus := 1
		xConnectionIndex := connectionIndex{}
		w := *response.NewWriter(httptest.NewRecorder())

		registerConnection(w, xConnectionIndex, key)
		_, values := collections.MapEntries(xConnectionIndex[key])
		assert.HasLength(t, values, 1)
		values[0].Status = xStatus
	})

	t.Run("MultipleConnectionsOnKey", func(t *testing.T) {
		var key int64 = 1
		w := *response.NewWriter(httptest.NewRecorder())
		xConnectionIndex := connectionIndex{}

		registerConnection(w, xConnectionIndex, key)
		registerConnection(w, xConnectionIndex, key)
		_, values := collections.MapEntries(xConnectionIndex[key])
		assert.HasLength(t, values, 2)
	})
}

func TestCloseSSEConnection(t *testing.T) {
	setup := func() connectionIndex {
		w := *response.NewWriter(httptest.NewRecorder())
		return connectionIndex{
			0: {"a": w},
			1: {"a": w, "b": w, "c": w},
		}
	}

	t.Run("Normal", func(t *testing.T) {
		xConnectionIndex := setup()
		var controlKey int64 = 0
		var testKey int64 = 1

		closeSSEConnection(xConnectionIndex, testKey, "b")
		keys, _ := collections.MapEntries(xConnectionIndex[controlKey])
		assert.HasLength(t, keys, 1)
		keys, _ = collections.MapEntries(xConnectionIndex[testKey])
		sort.Strings(keys)
		assert.HasLength(t, keys, 2)
		assert.Equals(t, keys[0], "a")
		assert.Equals(t, keys[1], "c")
	})

	t.Run("CleansEmptyKeys", func(t *testing.T) {
		xConnectionIndex := setup()
		var emptyKey int64 = 0
		var nonEmptyKey int64 = 1

		closeSSEConnection(xConnectionIndex, emptyKey, "a")
		keys, _ := collections.MapEntries(xConnectionIndex[nonEmptyKey])
		assert.HasLength(t, keys, 3)
		_, ok := xConnectionIndex[emptyKey]
		assert.Equals(t, ok, false)
	})
}
