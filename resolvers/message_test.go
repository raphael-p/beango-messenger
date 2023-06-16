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

var chatID int = 1

func setupMessageTests(t *testing.T, body string) (
	*response.Writer,
	*http.Request,
	database.Connection,
) {
	w, req := mockRequest(body)
	conn := mocks.MakeMockConnection()
	param := map[string]string{CHAT_ID_KEY: fmt.Sprint(chatID)}
	req = setContext(t, req, mocks.Admin, param)
	return w, req, conn
}

func TestGetChatMessages(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w, req, conn := setupMessageTests(t, "")
		userID1 := mocks.ADMIN_ID
		userID2 := 12
		conn.SetChat(mocks.MakePrivateChat(chatID), userID1, userID2)
		conn.SetMessage(mocks.MakeMessage(1, userID1, chatID))
		conn.SetMessage(mocks.MakeMessage(2, userID2, chatID))

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		messages := &[]database.Message{}
		err := json.Unmarshal([]byte(w.Body), messages)
		assert.IsNil(t, err)
		assert.HasLength(t, *messages, 2)
	})

	t.Run("NoMessages", func(t *testing.T) {
		w, req, conn := setupMessageTests(t, "")
		conn.SetChat(mocks.MakePrivateChat(chatID), mocks.Admin.ID, 11)

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Equals(t, w.Body, "[]")
	})

	t.Run("NoChat", func(t *testing.T) {
		w, req, conn := setupMessageTests(t, "")

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "chat not found")
	})

	t.Run("NotChatUser", func(t *testing.T) {
		w, req, conn := setupMessageTests(t, "")
		conn.SetChat(mocks.MakePrivateChat(chatID), 11, 12)

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "chat not found")
	})
}

func TestSendMessage(t *testing.T) {
	content := "Hello, World!"
	body := fmt.Sprintf(`{"content": "%s"}`, content)

	t.Run("Normal", func(t *testing.T) {
		w, req, conn := setupMessageTests(t, body)
		conn.SetChat(mocks.MakePrivateChat(chatID), mocks.Admin.ID, 12)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusAccepted)
		message := &database.Message{}
		err := json.Unmarshal([]byte(w.Body), message)
		assert.IsNil(t, err)
		assert.Equals(t, message.UserID, mocks.Admin.ID)
		assert.Equals(t, message.ChatID, chatID)
		assert.Equals(t, message.Content, content)
	})

	t.Run("NoChat", func(t *testing.T) {
		w, req, conn := setupMessageTests(t, body)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "chat not found")
	})

	t.Run("NotChatUser", func(t *testing.T) {
		w, req, conn := setupMessageTests(t, body)
		conn.SetChat(mocks.MakePrivateChat(chatID), 11, 12)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "chat not found")
	})
}
