package resolvers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/response"
)

func setupChatTests(t *testing.T, body string) (
	*response.Writer,
	*http.Request,
	database.Connection,
) {
	w, req := mockRequest(body)
	conn := mocks.MakeMockConnection()
	req = setContext(t, req, mocks.Admin, nil)
	return w, req, conn
}

func TestGetChats(t *testing.T) {
	adminID := mocks.ADMIN_ID
	uuid := func() string { return uuid.NewString() }
	testCases := []struct {
		name          string
		chatUsers     [][2]string
		expectedCount int
	}{
		{
			"UserInAllChats",
			[][2]string{{adminID, uuid()}, {uuid(), adminID}},
			2,
		},
		{
			"UserInSomeChats",
			[][2]string{{adminID, uuid()}, {uuid(), uuid()}},
			1,
		},
		{
			"UserInSomeChats",
			[][2]string{{uuid(), uuid()}, {uuid(), uuid()}},
			0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			w, req, conn := setupChatTests(t, "")
			for _, pair := range testCase.chatUsers {
				conn.SetChat(mocks.MakeChat(pair[0], pair[1]))
			}

			GetChats(w, req, conn)
			assert.Equals(t, w.Status, http.StatusOK)
			chats := &[]database.Chat{}
			err := json.Unmarshal([]byte(w.Body), chats)
			assert.IsNil(t, err)
			assert.HasLength(t, *chats, testCase.expectedCount)
		})
	}
}

func TestCreateChat(t *testing.T) {
	body := func(userID string) string {
		return fmt.Sprintf(`{"userID": "%s"}`, userID)
	}

	createAndCheck := func(
		w *response.Writer,
		req *http.Request,
		conn database.Connection,
		userID string,
	) {
		CreateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusCreated)
		chat := &database.Chat{}
		err := json.Unmarshal([]byte(w.Body), chat)
		assert.IsNil(t, err)
		xUserIDs := [2]string{mocks.ADMIN_ID, userID}
		assert.Equals(t, chat.UserIDs, xUserIDs)
	}

	t.Run("Normal", func(t *testing.T) {
		user := mocks.MakeUser()
		w, req, conn := setupChatTests(t, body(user.ID))
		conn.SetUser(user)
		createAndCheck(w, req, conn, user.ID)
	})

	t.Run("SelfChat", func(t *testing.T) {
		w, req, conn := setupChatTests(t, body(mocks.ADMIN_ID))
		createAndCheck(w, req, conn, mocks.ADMIN_ID)
	})

	t.Run("UserDoesNotExist", func(t *testing.T) {
		user := mocks.MakeUser()
		w, req, conn := setupChatTests(t, body(user.ID))

		CreateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		xError := fmt.Sprintf("userID %s is invalid", user.ID)
		assert.Equals(t, w.Body, xError)
	})

	t.Run("ChatAlreadyExists", func(t *testing.T) {
		user := mocks.MakeUser()
		w, req, conn := setupChatTests(t, body(user.ID))
		conn.SetUser(user)
		chat := mocks.MakeChat(mocks.ADMIN_ID, user.ID)
		conn.SetChat(chat)

		CreateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusConflict)
		assert.Equals(t, w.Body, "chat already exists")
	})
}
