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

func TestGetChatMessages(t *testing.T) {
	setup := func(contextUser *database.User) (
		*response.Writer,
		*http.Request,
		database.Connection,
		*database.Chat,
	) {
		w, req := mockRequest("")
		conn := mocks.MakeMockConnection()
		chat := mocks.MakeChat(mocks.Admin.ID, mocks.MakeUser().ID)
		param := map[string]string{"chatID": chat.ID}
		if contextUser == nil {
			contextUser = mocks.Admin
		}
		req = setContext(t, req, contextUser, param)
		return w, req, conn, chat
	}

	t.Run("Normal", func(t *testing.T) {
		w, req, conn, chat := setup(nil)
		conn.SetChat(chat)
		conn.SetMessage(mocks.MakeMessage(chat.UserIDs[0], chat.ID))
		conn.SetMessage(mocks.MakeMessage(chat.UserIDs[1], chat.ID))

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		messages := &[]database.Message{}
		err := json.Unmarshal([]byte(w.Body), messages)
		assert.IsNil(t, err)
		assert.HasLength(t, *messages, 2)
	})

	t.Run("NoMessages", func(t *testing.T) {
		w, req, conn, chat := setup(nil)
		conn.SetChat(chat)

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Equals(t, w.Body, "[]")
	})

	t.Run("NoChat", func(t *testing.T) {
		w, req, conn, _ := setup(nil)

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "chat not found")
	})

	t.Run("NotChatUser", func(t *testing.T) {
		w, req, conn, chat := setup(mocks.MakeUser())
		conn.SetChat(chat)

		GetChatMessages(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "chat not found")
	})
}

func TestSendMessage(t *testing.T) {
	setup := func(content string, contextUser *database.User) (
		*response.Writer,
		*http.Request,
		database.Connection,
		*database.Chat,
	) {
		body := fmt.Sprintf(`{"content": "%s"}`, content)
		w, req := mockRequest(body)
		conn := mocks.MakeMockConnection()
		chat := mocks.MakeChat(mocks.Admin.ID, mocks.MakeUser().ID)
		param := map[string]string{"chatID": chat.ID}
		if contextUser == nil {
			contextUser = mocks.Admin
		}
		req = setContext(t, req, contextUser, param)
		return w, req, conn, chat
	}

	t.Run("Normal", func(t *testing.T) {
		content := "Hello, World!"
		w, req, conn, chat := setup(content, nil)
		conn.SetChat(chat)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusAccepted)
		message := &database.Message{}
		err := json.Unmarshal([]byte(w.Body), message)
		assert.IsNil(t, err)
		assert.Equals(t, message.UserID, mocks.Admin.ID)
		assert.Equals(t, message.ChatID, chat.ID)
		assert.Equals(t, message.Content, content)
	})

	t.Run("NoChat", func(t *testing.T) {
		content := "Hello, World!"
		w, req, conn, _ := setup(content, nil)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "chat not found")
	})

	t.Run("NotChatUser", func(t *testing.T) {
		content := "Hello, World!"
		w, req, conn, chat := setup(content, mocks.MakeUser())
		conn.SetChat(chat)

		SendMessage(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "chat not found")
	})
}
