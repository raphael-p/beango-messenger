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
	testCases := []struct {
		name          string
		chatUsers     [2][]int
		expectedCount int
	}{
		{
			"UserInAllChats",
			[2][]int{{adminID, 11}, {12, adminID}},
			2,
		},
		{
			"UserInSomeChats",
			[2][]int{{adminID, 13}, {14, 15}},
			1,
		},
		{
			"UserInNoChats",
			[2][]int{{16, 17}, {18, 19}},
			0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			w, req, conn := setupChatTests(t, "")
			for _, pair := range testCase.chatUsers {
				conn.SetChat(mocks.MakePrivateChat(1), pair...)
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

func TestCreatePrivateChat(t *testing.T) {
	body := func(userID int) string {
		return fmt.Sprintf(`{"userID": %d}`, userID)
	}

	createAndCheck := func(
		w *response.Writer,
		req *http.Request,
		conn database.Connection,
	) {
		CreatePrivateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusCreated)
		chat := &database.Chat{}
		err := json.Unmarshal([]byte(w.Body), chat)
		assert.IsNil(t, err)
		assert.Equals(t, chat.Name, "")
		assert.Equals(t, chat.ChatType, database.PRIVATE_CHAT)
	}

	t.Run("Normal", func(t *testing.T) {
		user := mocks.MakeUser(11)
		w, req, conn := setupChatTests(t, body(user.ID))
		conn.SetUser(user)
		createAndCheck(w, req, conn)
	})

	t.Run("SelfChat", func(t *testing.T) {
		w, req, conn := setupChatTests(t, body(mocks.ADMIN_ID))
		createAndCheck(w, req, conn)
	})

	t.Run("UserDoesNotExist", func(t *testing.T) {
		user := mocks.MakeUser(11)
		w, req, conn := setupChatTests(t, body(user.ID))

		CreatePrivateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		xError := fmt.Sprintf("userID %d is invalid", user.ID)
		assert.Equals(t, w.Body, xError)
	})

	t.Run("ChatAlreadyExists", func(t *testing.T) {
		user := mocks.MakeUser(11)
		w, req, conn := setupChatTests(t, body(user.ID))
		conn.SetUser(user)
		chat := mocks.MakePrivateChat(1)
		conn.SetChat(chat, mocks.ADMIN_ID, user.ID)

		CreatePrivateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusConflict)
		assert.Equals(t, w.Body, "chat already exists")
	})
}
