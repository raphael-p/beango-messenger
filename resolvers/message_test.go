package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/response"
)

func setupMessageTests(t *testing.T, body string, userID1, userID2 int64) (
	*response.Writer,
	*http.Request,
	database.Connection,
	int64,
) {
	w, req := mockRequest(body)
	conn := mocks.MakeMockConnection()
	var chatID int64
	if userID1 != 0 && userID2 != 0 {
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), userID1, userID2)
		chatID = chat.ID
	}
	param := map[string]string{CHAT_ID_KEY: fmt.Sprint(chatID)}
	req = setContext(t, req, mocks.Admin, param)
	return w, req, conn, chatID
}

func TestGetChatMessages(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		userID1 := mocks.ADMIN_ID
		var userID2 int64 = 12
		w, req, conn, chatID := setupMessageTests(t, "", userID1, userID2)
		conn.SetMessage(mocks.MakeMessage(userID1, chatID))
		conn.SetMessage(mocks.MakeMessage(userID2, chatID))

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		messages := &[]database.MessageDatabase{}
		err := json.Unmarshal(w.Body, messages)
		assert.IsNil(t, err)
		assert.HasLength(t, *messages, 2)
	})

	t.Run("NoMessages", func(t *testing.T) {
		w, req, conn, _ := setupMessageTests(t, "", mocks.ADMIN_ID, 11)

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Equals(t, string(w.Body), "[]")
	})

	t.Run("NoChat", func(t *testing.T) {
		w, req, conn, _ := setupMessageTests(t, "", 0, 0)

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, string(w.Body), "chat not found")
	})

	t.Run("NotChatUser", func(t *testing.T) {
		w, req, conn, _ := setupMessageTests(t, "", 11, 12)

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, string(w.Body), "chat not found")
	})
}

func TestSendMessage(t *testing.T) {
	content := "Hello, World!"
	body := fmt.Sprintf(`{"content": "%s"}`, content)

	t.Run("Normal", func(t *testing.T) {
		w, req, conn, chatID := setupMessageTests(t, body, mocks.ADMIN_ID, 12)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusAccepted)
		message := &database.MessageDatabase{}
		err := json.Unmarshal(w.Body, message)
		assert.IsNil(t, err)
		assert.Equals(t, message.UserID, mocks.Admin.ID)
		assert.Equals(t, message.ChatID, chatID)
		assert.Equals(t, message.Content, content)
	})

	t.Run("NoChat", func(t *testing.T) {
		w, req, conn, _ := setupMessageTests(t, body, 0, 0)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, string(w.Body), "chat not found")
	})

	t.Run("NotChatUser", func(t *testing.T) {
		w, req, conn, _ := setupMessageTests(t, body, 11, 12)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, string(w.Body), "chat not found")
	})
}
