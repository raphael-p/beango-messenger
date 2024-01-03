package resolvers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/collections"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func TestRegisterChatSSE(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		buf := logger.MockFileLogger(t)
		w, r, conn := resolverutils.CommonSetup("")
		user, _ := conn.SetUser(mocks.MakeUser())
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), user.ID, mocks.ADMIN_ID)
		params := map[string]string{resolverutils.CHAT_ID_KEY: fmt.Sprint(chat.ID)}
		r = resolverutils.SetContext(t, r, mocks.Admin, params)

		done := make(chan bool)
		go func() {
			RegisterChatSSE(w, r, conn)
			done <- true
		}()

		select {
		case <-done:
			t.Error("connection closed unexpectedly")
		case <-time.After(1 * time.Second):
			assert.Contains(t, buf.String(), "opened")
		}
	})

	t.Run("RedirectsOnMissingUser", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		params := map[string]string{resolverutils.CHAT_ID_KEY: "1"}
		r = resolverutils.SetContext(t, r, nil, params)

		RegisterChatSSE(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), "redirect")
	})

	t.Run("RedirectsOnMissingChatID", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		r = resolverutils.SetContext(t, r, mocks.Admin, nil)

		RegisterChatSSE(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), "redirect")
	})
}

func TestSendChatEvent(t *testing.T) {
	setup := func() (*bytes.Buffer, int64, string) {
		buf := logger.MockFileLogger(t)
		var key int64 = 1
		connectionID := "a"
		w := response.NewWriter(httptest.NewRecorder())
		chatConnectionIndex[key] = connectionMap{connectionID: w}
		return buf, key, connectionID
	}

	t.Run("Normal", func(t *testing.T) {
		buf, key, connectionID := setup()
		xEvent := "test-event"
		xMessage := fmt.Sprintf("[SSE connection %s] sent '%s' event", connectionID, xEvent)

		SendChatEvent(key, xEvent, "Hello World!")
		assert.Contains(t, buf.String(), xMessage)
	})

	t.Run("ChatNotFound", func(t *testing.T) {
		buf, _, _ := setup()

		SendChatEvent(2, "test-event", "Hello World!")
		assert.Equals(t, buf.String(), "")
	})
}

func TestRegisterConnection(t *testing.T) {
	setup := func() (*response.Writer, connectionIndex, int64) {
		var key int64 = 1
		w := response.NewWriter(httptest.NewRecorder())
		return w, connectionIndex{}, key
	}

	t.Run("Normal", func(t *testing.T) {
		w, xConnectionIndex, key := setup()
		xStatus := 1

		registerConnection(w, xConnectionIndex, key)
		_, values := collections.MapEntries(xConnectionIndex[key])
		assert.HasLength(t, values, 1)
		values[0].Status = xStatus
	})

	t.Run("MultipleConnectionsOnKey", func(t *testing.T) {
		w, xConnectionIndex, key := setup()

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
