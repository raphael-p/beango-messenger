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

func TestGenerateChatName(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		users := []database.User{
			{ID: 98, DisplayName: "TARIK"},
			{ID: 11, DisplayName: "GEORGE"},
			{ID: 12, DisplayName: "AMY"},
		}
		chatName := generateChatName(98, users)
		assert.Equals(t, chatName, "AMY, GEORGE")
	})

	t.Run("OrderAgnostic", func(t *testing.T) {
		users := []database.User{
			{ID: 32, DisplayName: "TARIK"},
			{ID: 5, DisplayName: "GEORGE"},
			{ID: 34, DisplayName: "AMY"},
		}
		chatName := generateChatName(5, users)
		assert.Equals(t, chatName, "AMY, TARIK")

		users = []database.User{
			{ID: 34, DisplayName: "AMY"},
			{ID: 5, DisplayName: "GEORGE"},
			{ID: 32, DisplayName: "TARIK"},
		}
		chatName = generateChatName(5, users)
		assert.Equals(t, chatName, "AMY, TARIK")
	})

	t.Run("NoOtherUsers", func(t *testing.T) {
		users := []database.User{
			{ID: 5, DisplayName: "GEORGE"},
		}
		chatName := generateChatName(5, users)
		assert.Equals(t, chatName, "")
	})

	t.Run("EmptyUserSlice", func(t *testing.T) {
		users := []database.User{}
		chatName := generateChatName(5, users)
		assert.Equals(t, chatName, "")
	})

	t.Run("UserNotInSlice", func(t *testing.T) {
		users := []database.User{
			{ID: 98, DisplayName: "TARIK"},
			{ID: 11, DisplayName: "GEORGE"},
			{ID: 12, DisplayName: "AMY"},
		}
		chatName := generateChatName(1, users)
		assert.Equals(t, chatName, "AMY, GEORGE, TARIK")
	})
}

func TestChatsDatabase(t *testing.T) {
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
			for _, pair := range testCase.chatUsers {
				conn.SetChat(mocks.MakePrivateChat(), pair...)
			}

			chats, httpError := chatsDatabase(adminID, conn)
			assert.IsNil(t, httpError)
			assert.HasLength(t, chats, testCase.expectedCount)
		})
	}
}

func TestGetChats(t *testing.T) {
	adminID := mocks.ADMIN_ID

	t.Run("Normal", func(t *testing.T) {
		w, req, conn := resolverutils.CommonSetup("")
		req = resolverutils.SetContext(t, req, mocks.Admin, nil)
		conn.SetChat(mocks.MakePrivateChat(), adminID, 13)

		GetChats(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		chats := &[]database.Chat{}
		err := json.Unmarshal(w.Body, chats)
		assert.IsNil(t, err)
		assert.HasLength(t, *chats, 1)
	})
}

func TestCreatePrivateChat(t *testing.T) {
	setup := func(t *testing.T, userID int64) (*response.Writer, *http.Request) {
		w, req, _ := resolverutils.CommonSetup(fmt.Sprintf(`{"userID": %d}`, userID))
		req = resolverutils.SetContext(t, req, mocks.Admin, nil)
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

		CreatePrivateChat(w, req, conn)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Equals(t, string(w.Body), "cannot create a chat with yourself")
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
