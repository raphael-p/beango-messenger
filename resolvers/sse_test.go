package resolvers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/utils/collections"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func TestSendChatEvent(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		var key int64 = 1
		connectionID := "a"
		w := response.NewWriter(httptest.NewRecorder())
		chatConnectionIndex[key] = connectionMap{connectionID: w}
		buf := logger.MockFileLogger(t)
		xEvent := "test-event"
		xMessage := fmt.Sprintf("[SSE connection %s] sent '%s' event", connectionID, xEvent)

		SendChatEvent(key, xEvent, "Hello World!")
		assert.Contains(t, buf.String(), xMessage)
	})

	t.Run("ChatNotFound", func(t *testing.T) {
		var key int64 = 1
		connectionID := "a"
		w := response.NewWriter(httptest.NewRecorder())
		chatConnectionIndex[key] = connectionMap{connectionID: w}
		buf := logger.MockFileLogger(t)

		SendChatEvent(2, "test-event", "Hello World!")
		assert.Equals(t, buf.String(), "")
	})
}

func TestRegisterConnection(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		var key int64 = 1
		xStatus := 1
		xConnectionIndex := connectionIndex{}
		w := response.NewWriter(httptest.NewRecorder())

		registerConnection(w, xConnectionIndex, key)
		_, values := collections.MapEntries(xConnectionIndex[key])
		assert.HasLength(t, values, 1)
		values[0].Status = xStatus
	})

	t.Run("MultipleConnectionsOnKey", func(t *testing.T) {
		var key int64 = 1
		w := response.NewWriter(httptest.NewRecorder())
		xConnectionIndex := connectionIndex{}

		registerConnection(w, xConnectionIndex, key)
		registerConnection(w, xConnectionIndex, key)
		_, values := collections.MapEntries(xConnectionIndex[key])
		assert.HasLength(t, values, 2)
	})
}

func TestTrapConnection(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		var key int64 = 1
		connectionID := "a"
		w := response.NewWriter(httptest.NewRecorder())
		connectionIndex := connectionIndex{key: {connectionID: w}}
		buf := logger.MockFileLogger(t)

		ctx, cancel := context.WithCancel(context.Background())
		r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/test", nil)
		done := make(chan bool)
		go func() {
			trapConnection(r, connectionIndex, key, connectionID)
			done <- true
		}()

		time.Sleep(100 * time.Millisecond)
		assert.Contains(t, buf.String(), "opened")
		assert.NotContains(t, buf.String(), "closed")

		cancel()
		select {
		case <-done:
			assert.Contains(t, buf.String(), "closed")
		case <-time.After(1 * time.Second):
			t.Error("test timed out")
		}
	})
}

func TestCloseSSEConnection(t *testing.T) {
	setup := func() connectionIndex {
		w := response.NewWriter(httptest.NewRecorder())
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

func TestSendEvent(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w := response.NewWriter(httptest.NewRecorder())
		connectionID := "a"
		connections := connectionMap{connectionID: w}
		buf := logger.MockFileLogger(t)
		xEvent := "test-event"
		xData := "Hello World!"
		xMessage := fmt.Sprintf("[SSE connection %s] sent '%s' event", connectionID, xEvent)

		sendEvent(connections, xEvent, xData)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), xEvent, xData)
		assert.Contains(t, buf.String(), xMessage)
	})

	t.Run("SendsMultiple", func(t *testing.T) {
		w1 := response.NewWriter(httptest.NewRecorder())
		w2 := response.NewWriter(httptest.NewRecorder())
		connections := connectionMap{"a": w1, "b": w2}
		xEvent := "test-event"
		xData := "Hello World!"

		sendEvent(connections, xEvent, xData)
		assert.Equals(t, w1.Status, http.StatusOK)
		assert.Contains(t, string(w1.Body), xEvent, xData)
		assert.Equals(t, w2.Status, http.StatusOK)
		assert.Contains(t, string(w2.Body), xEvent, xData)
	})
}
