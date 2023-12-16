package resolvers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

var sseConnections = map[int64]map[string]response.Writer{}

func RegisterSSE(w *response.Writer, r *http.Request, conn database.Connection) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// NOTE: uncomment next 2 lines if doesn't work in production
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Status = 200 // avoids WriteHeader calls
	newWriter := *w

	user, _, httpError := resolverutils.GetRequestContext(r)
	if httpError != nil {
		w.WriteSSE("redirect", "connection expired")
		return
	}

	sseConnectionID := uuid.NewString()
	connectionForUser := sseConnections[user.ID]
	if connectionForUser == nil {
		connectionForUser = map[string]response.Writer{sseConnectionID: newWriter}
	} else {
		connectionForUser[sseConnectionID] = newWriter
	}
	sseConnections[user.ID] = connectionForUser

	dataCh := make(chan string) // DEBUG

	// Cleanup
	_, cancel := context.WithCancel(r.Context())
	defer func() {
		cancel()
		close(dataCh) // DEBUG
		dataCh = nil  // DEBUG
		closeSSEConnection(user.ID, sseConnectionID)
		message := fmt.Sprintf("SSE connection %s closed for user %d", sseConnectionID, user.ID)
		logger.Info(message)
	}()

	// DEBUG
	go func() {
		for data := range dataCh {
			div := fmt.Sprintf(`<div id="errors" hx-swap-oob="innerHTML">%s</div>`, data)
			SendEvent(user.ID, "message", div)
		}
	}()

	message := fmt.Sprintf("SSE connection %s opened for user %d", sseConnectionID, user.ID)
	logger.Info(message)
	for {
		select {
		case <-r.Context().Done():
			return
		default:
			// check for reasons to break the connection
			if _, ok := sseConnections[user.ID][sseConnectionID]; !ok {
				return
			}
			sessionID, err := cookies.Get(r, cookies.SESSION)
			if err != nil {
				SendEvent(user.ID, "redirect", "")
				return
			}
			if _, ok := conn.CheckSession(sessionID); !ok {
				SendEvent(user.ID, "redirect", "")
				return
			}

			dataCh <- time.Now().Format(time.TimeOnly) // DEBUG
		}
		time.Sleep(time.Millisecond * 2000)
	}
}

func closeSSEConnection(userID int64, connectionID string) {
	delete(sseConnections[userID], connectionID)
	if len(sseConnections[userID]) == 0 {
		delete(sseConnections, userID)
	}
}

// Writes to all SSE connections for a given user.
// A user will have multiple connections open if they
// open multiple tabs, for instance.
func SendEvent(userID int64, event, data string) {
	connections := sseConnections[userID]
	for connectionID, w := range connections {
		w.WriteSSE(event, data)
		message := fmt.Sprintf(
			"sent '%s' event to user %d on connection %s",
			event,
			userID,
			connectionID,
		)
		logger.Info(message)
	}
}
