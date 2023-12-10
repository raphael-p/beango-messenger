package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/response"
)

func setupMessageTests(userID1, userID2 int64) (
	database.Connection,
	int64,
) {
	conn := mocks.MakeMockConnection()
	var chatID int64
	if userID1 != 0 && userID2 != 0 {
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), userID1, userID2)
		chatID = chat.ID
	}
	return conn, chatID
}

func makeMessageRequest(t *testing.T, body string, chatID int64) (*response.Writer, *http.Request) {
	w, req, _ := resolverutils.CommonSetup(body)
	params := map[string]string{resolverutils.CHAT_ID_KEY: fmt.Sprint(chatID)}
	req = resolverutils.SetContext(t, req, mocks.Admin, params)
	return w, req
}

func TestChatMessagesDatabase(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		userID1 := mocks.ADMIN_ID
		var userID2 int64 = 12
		conn, chatID := setupMessageTests(userID1, userID2)
		conn.SetMessage(mocks.MakeMessage(userID1, chatID))
		conn.SetMessage(mocks.MakeMessage(userID2, chatID))

		messages, httpError := chatMessagesDatabase(mocks.ADMIN_ID, chatID, 0, 0, 0, conn)
		assert.IsNil(t, httpError)
		assert.HasLength(t, messages, 2)
	})

	t.Run("NoMessages", func(t *testing.T) {
		conn, chatID := setupMessageTests(mocks.ADMIN_ID, 11)

		messages, httpError := chatMessagesDatabase(mocks.ADMIN_ID, chatID, 0, 0, 0, conn)
		assert.IsNil(t, httpError)
		assert.HasLength(t, messages, 0)
	})

	t.Run("NoChat", func(t *testing.T) {
		conn, chatID := setupMessageTests(0, 0)

		messages, httpError := chatMessagesDatabase(mocks.ADMIN_ID, chatID, 0, 0, 0, conn)
		assert.IsNil(t, messages)
		resolverutils.AssertHTTPError(t, httpError, http.StatusNotFound, "chat not found")
	})

	t.Run("NotChatUser", func(t *testing.T) {
		conn, chatID := setupMessageTests(11, 12)
		messages, httpError := chatMessagesDatabase(mocks.ADMIN_ID, chatID, 0, 0, 0, conn)
		assert.IsNil(t, messages)
		resolverutils.AssertHTTPError(t, httpError, http.StatusNotFound, "chat not found")
	})
}

func TestGetChatMessages(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		userID1 := mocks.ADMIN_ID
		var userID2 int64 = 12

		conn, chatID := setupMessageTests(userID1, userID2)
		w, req := makeMessageRequest(t, "", chatID)
		conn.SetMessage(mocks.MakeMessage(userID1, chatID))

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		messages := &[]database.MessageDatabase{}
		err := json.Unmarshal(w.Body, messages)
		assert.IsNil(t, err)
		assert.HasLength(t, *messages, 1)
	})
}

func TestSendMessage(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		content := "Hello, World!"
		body := fmt.Sprintf(`{"content": "%s"}`, content)
		conn, chatID := setupMessageTests(mocks.ADMIN_ID, 12)
		w, req := makeMessageRequest(t, body, chatID)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusCreated)
		message := &database.MessageDatabase{}
		err := json.Unmarshal(w.Body, message)
		assert.IsNil(t, err)
		assert.Equals(t, message.UserID, mocks.Admin.ID)
		assert.Equals(t, message.ChatID, chatID)
		assert.Equals(t, message.Content, content)
	})
}

func TestSendMessageDatabase(t *testing.T) {
	content := "Hello, World!"

	t.Run("Normal", func(t *testing.T) {
		conn, chatID := setupMessageTests(mocks.ADMIN_ID, 12)
		message, httpError := sendMessageDatabase(mocks.ADMIN_ID, chatID, content, conn)
		assert.IsNil(t, httpError)
		assert.Equals(t, message.UserID, mocks.Admin.ID)
		assert.Equals(t, message.ChatID, chatID)
		assert.Equals(t, message.Content, content)
	})

	t.Run("NoChat", func(t *testing.T) {
		conn, chatID := setupMessageTests(0, 0)

		message, httpError := sendMessageDatabase(mocks.ADMIN_ID, chatID, content, conn)
		assert.IsNil(t, message)
		resolverutils.AssertHTTPError(t, httpError, http.StatusNotFound, "chat not found")
	})

	t.Run("NotChatUser", func(t *testing.T) {
		conn, chatID := setupMessageTests(11, 12)

		message, httpError := sendMessageDatabase(mocks.ADMIN_ID, chatID, content, conn)
		assert.IsNil(t, message)
		resolverutils.AssertHTTPError(t, httpError, http.StatusNotFound, "chat not found")
	})

	t.Run("TrimsSpace", func(t *testing.T) {
		paddedContent := " \n \r " + content + " \n \r "
		conn, chatID := setupMessageTests(mocks.ADMIN_ID, 12)
		message, httpError := sendMessageDatabase(mocks.ADMIN_ID, chatID, paddedContent, conn)
		assert.IsNil(t, httpError)
		assert.Equals(t, message.UserID, mocks.Admin.ID)
		assert.Equals(t, message.ChatID, chatID)
		assert.Equals(t, message.Content, content)
	})
}
