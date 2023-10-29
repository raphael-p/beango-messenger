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

func TestGetChatsDatabase(t *testing.T) {
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

			chats, httpError := getChatsDatabase(adminID, conn)
			assert.IsNil(t, httpError)
			assert.HasLength(t, chats, testCase.expectedCount)
		})
	}
}

func TestGetChats(t *testing.T) {
	adminID := mocks.ADMIN_ID

	t.Run("Normal", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		r = resolverutils.SetContext(t, r, mocks.Admin, nil)
		conn.SetChat(mocks.MakePrivateChat(), adminID, 13)

		GetChats(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		chats := &[]database.Chat{}
		err := json.Unmarshal(w.Body, chats)
		assert.IsNil(t, err)
		assert.HasLength(t, *chats, 1)
	})
}

func TestValidateCreatePrivateChatInput(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		input := createPrivateChatInput{123}

		httpError := validateCreatePrivateChatInput(&input, 1234)
		assert.IsNil(t, httpError)
	})

	t.Run("SelfChat", func(t *testing.T) {
		input := createPrivateChatInput{123}

		httpError := validateCreatePrivateChatInput(&input, 123)
		assert.Equals(t, httpError.Status, http.StatusBadRequest)
		assert.Equals(t, httpError.Message, "cannot create a chat with yourself")
	})
}

func TestCreatePrivateChatDatabase(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		user, _ := conn.SetUser(mocks.MakeUser())

		chat, httpError := createPrivateChatDatabase(mocks.ADMIN_ID, user.ID, conn)
		assert.IsNil(t, httpError)
		assert.Equals(t, chat.Name, "")
		assert.Equals(t, chat.Type, database.PRIVATE_CHAT)
	})

	t.Run("UserDoesNotExist", func(t *testing.T) {
		var fakeUserID int64 = 451
		conn := mocks.MakeMockConnection()

		chat, httpError := createPrivateChatDatabase(mocks.ADMIN_ID, fakeUserID, conn)
		assert.IsNil(t, chat)
		assert.Equals(t, httpError.Status, http.StatusBadRequest)
		xError := fmt.Sprintf("userID %d is invalid", fakeUserID)
		assert.Equals(t, httpError.Message, xError)
	})

	t.Run("ChatAlreadyExists", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		user, _ := conn.SetUser(mocks.MakeUser())
		chat := mocks.MakePrivateChat()
		conn.SetChat(chat, mocks.ADMIN_ID, user.ID)

		chat, httpError := createPrivateChatDatabase(mocks.ADMIN_ID, user.ID, conn)
		assert.IsNotNil(t, chat)
		assert.Equals(t, httpError.Status, http.StatusConflict)
		assert.Equals(t, httpError.Message, "chat already exists")
	})
}

func TestCreatePrivateChat(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		user, _ := conn.SetUser(mocks.MakeUser())
		w, r, _ := resolverutils.CommonSetup(fmt.Sprintf(`{"userID": %d}`, user.ID))
		r = resolverutils.SetContext(t, r, mocks.Admin, nil)

		CreatePrivateChat(w, r, conn)
		assert.Equals(t, w.Status, http.StatusCreated)
		assert.IsValidJSON(t, string(w.Body), &database.Chat{})
	})
}
