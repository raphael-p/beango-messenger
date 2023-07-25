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

func TestGetChats(t *testing.T) {
	adminID := mocks.ADMIN_ID
	testCases := []struct {
		name          string
		chatUsers     [2][]int64
		expectedCount int
	}{
		{
			"UserInAllChats",
			[2][]int64{{adminID, 11}, {12, adminID}},
			2,
		},
		{
			"UserInSomeChats",
			[2][]int64{{adminID, 13}, {14, 15}},
			1,
		},
		{
			"UserInNoChats",
			[2][]int64{{16, 17}, {18, 19}},
			0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			conn := mocks.MakeMockConnection()
			w, req := mockRequest("")
			req = setContext(t, req, mocks.Admin, nil)
			for _, pair := range testCase.chatUsers {
				conn.SetChat(mocks.MakePrivateChat(), pair...)
			}

			GetChats(w, req, conn)
			assert.Equals(t, w.Status, http.StatusOK)
			chats := &[]database.Chat{}
			err := json.Unmarshal(w.Body, chats)
			assert.IsNil(t, err)
			assert.HasLength(t, *chats, testCase.expectedCount)
		})
	}
}

func TestCreatePrivateChat(t *testing.T) {
	setup := func(t *testing.T, userID int64) (*response.Writer, *http.Request) {
		w, req := mockRequest(fmt.Sprintf(`{"userID": %d}`, userID))
		req = setContext(t, req, mocks.Admin, nil)
		return w, req
	}

	createAndCheck := func(
		w *response.Writer,
		req *http.Request,
		conn database.Connection,
	) {
		CreatePrivateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusCreated)
		chat := &database.Chat{}
		err := json.Unmarshal(w.Body, chat)
		assert.IsNil(t, err)
		assert.Equals(t, chat.Name, "")
		assert.Equals(t, chat.Type, database.PRIVATE_CHAT)
	}

	t.Run("Normal", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		user, _ := conn.SetUser(mocks.MakeUser())
		w, req := setup(t, user.ID)
		createAndCheck(w, req, conn)
	})

	t.Run("SelfChat", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		w, req := setup(t, mocks.Admin.ID)
		createAndCheck(w, req, conn)
	})

	t.Run("UserDoesNotExist", func(t *testing.T) {
		var fakeUserID int64 = 451
		conn := mocks.MakeMockConnection()
		w, req := setup(t, fakeUserID)

		CreatePrivateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		xError := fmt.Sprintf("userID %d is invalid", fakeUserID)
		assert.Equals(t, string(w.Body), xError)
	})

	t.Run("ChatAlreadyExists", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		user, _ := conn.SetUser(mocks.MakeUser())
		w, req := setup(t, user.ID)
		chat := mocks.MakePrivateChat()
		conn.SetChat(chat, mocks.ADMIN_ID, user.ID)

		CreatePrivateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusConflict)
		assert.Equals(t, string(w.Body), "chat already exists")
	})
}
