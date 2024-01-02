package resolvers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

type connectionMap = map[string]*response.Writer
type connectionIndex = map[int64]connectionMap

var chatConnectionIndex = connectionIndex{}

// TODO: test
func RegisterChatSSE(w *response.Writer, r *http.Request, conn database.Connection) {
	newWriter := upgradeConnection(w)

	_, params, httpError := resolverutils.GetRequestContext(r, resolverutils.CHAT_ID_KEY)
	if httpError != nil {
		w.WriteSSE("redirect", "")
		return
	}
	chatID := params.ChatID

	chatConnectionIndex, connectionID := registerConnection(newWriter, chatConnectionIndex, chatID)
	trapConnection(r, chatConnectionIndex, chatID, connectionID)
}

// TODO: test
func SendChatEvent(chatID int64, event, data string) {
	chatConnections, ok := chatConnectionIndex[chatID]
	if ok {
		sendEvent(chatConnections, event, data)
	}
}

// Upgrades an HTTP connection to an SSE connection
func upgradeConnection(w *response.Writer) *response.Writer {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// NOTE: uncomment next 2 lines if doesn't work in production
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Status = 200 // having this set avoids unnecessary WriteHeader calls
	return w
}

// Add connection to the index.
func registerConnection(
	w *response.Writer,
	index connectionIndex,
	key int64,
) (connectionIndex, string) {
	sseConnectionID := uuid.NewString()
	connectionForKey := index[key]
	if connectionForKey == nil {
		connectionForKey = connectionMap{sseConnectionID: w}
	} else {
		connectionForKey[sseConnectionID] = w
	}
	index[key] = connectionForKey
	return index, sseConnectionID
}

// Makes sure connection is kept alive until terminated by client
func trapConnection(r *http.Request, index connectionIndex, key int64, connectionID string) {
	ctx, cancel := context.WithCancel(r.Context())
	defer func() {
		cancel()
		closeSSEConnection(index, key, connectionID)
		message := fmt.Sprintf("[SSE connection %s] closed", connectionID)
		logger.Info(message)
	}()

	message := fmt.Sprintf("[SSE connection %s] opened", connectionID)
	logger.Info(message)
	<-ctx.Done() // wait for client termination
}

// Removes an SSE connection from the index. Will remove an index entry if there
// are no more connections against it.
func closeSSEConnection(index connectionIndex, key int64, connectionID string) {
	delete(index[key], connectionID)
	if len(index[key]) == 0 {
		delete(index, key)
	}
}

// Writes to all SSE connections for a given user. A user will
// have multiple connections open if they open multiple tabs,
// for instance.
func sendEvent(connections connectionMap, event, data string) {
	for connectionID, w := range connections {
		w.WriteSSE(event, data)
		message := fmt.Sprintf(
			"[SSE connection %s] sent '%s' event",
			connectionID,
			event,
		)
		logger.Info(message)
	}
}
