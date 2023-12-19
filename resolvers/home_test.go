package resolvers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
)

func TestHome(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		r = resolverutils.SetContext(t, r, mocks.Admin, nil)

		Home(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), "<html>", "</html")
	})
}

func TestOpenChat(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		user, _ := conn.SetUser(mocks.MakeUser())
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), user.ID, mocks.ADMIN_ID)
		params := map[string]string{resolverutils.CHAT_ID_KEY: fmt.Sprint(chat.ID)}
		r = resolverutils.SetContext(t, r, mocks.Admin, params)
		query := r.URL.Query()
		query.Add("name", "My Chat Name")
		r.URL.RawQuery = query.Encode()

		OpenChat(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), "<table", "</table>")
	})
}

func TestRefreshMessages(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		user, _ := conn.SetUser(mocks.MakeUser())
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), user.ID, mocks.ADMIN_ID)
		message, _ := conn.SetMessage(mocks.MakeMessage(user.ID, chat.ID))
		params := map[string]string{resolverutils.CHAT_ID_KEY: fmt.Sprint(chat.ID)}
		r = resolverutils.SetContext(t, r, mocks.Admin, params)
		query := r.URL.Query()
		query.Add("from", "0")
		r.URL.RawQuery = query.Encode()

		RefreshMessages(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), message.Content, "<table", "</table>")
	})

	t.Run("NoNewMessages", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		user, _ := conn.SetUser(mocks.MakeUser())
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), user.ID, mocks.ADMIN_ID)
		params := map[string]string{resolverutils.CHAT_ID_KEY: fmt.Sprint(chat.ID)}
		r = resolverutils.SetContext(t, r, mocks.Admin, params)
		query := r.URL.Query()
		query.Add("from", "0")
		r.URL.RawQuery = query.Encode()

		RefreshMessages(w, r, conn)
		assert.Equals(t, w.Status, http.StatusNoContent)
		assert.Equals(t, string(w.Body), "")
	})
}

func TestScrollUp(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		user, _ := conn.SetUser(mocks.MakeUser())
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), user.ID, mocks.ADMIN_ID)
		message, _ := conn.SetMessage(mocks.MakeMessage(user.ID, chat.ID))
		params := map[string]string{resolverutils.CHAT_ID_KEY: fmt.Sprint(chat.ID)}
		r = resolverutils.SetContext(t, r, mocks.Admin, params)
		query := r.URL.Query()
		query.Add("to", "1")
		r.URL.RawQuery = query.Encode()

		ScrollUp(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), message.Content, "<table", "</table>")
	})

	t.Run("NoOlderMessages", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup("")
		user, _ := conn.SetUser(mocks.MakeUser())
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), user.ID, mocks.ADMIN_ID)
		params := map[string]string{resolverutils.CHAT_ID_KEY: fmt.Sprint(chat.ID)}
		r = resolverutils.SetContext(t, r, mocks.Admin, params)
		query := r.URL.Query()
		query.Add("to", "0")
		r.URL.RawQuery = query.Encode()

		ScrollUp(w, r, conn)
		assert.Equals(t, w.Status, http.StatusNoContent)
		assert.Equals(t, string(w.Body), "")
	})
}

func TestGetMessages(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		conn := mocks.MakeMockConnection()
		user, _ := conn.SetUser(mocks.MakeUser())
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), user.ID, mocks.ADMIN_ID)
		firstMessage, _ := conn.SetMessage(mocks.MakeMessage(user.ID, chat.ID))
		middleMessage, _ := conn.SetMessage(mocks.MakeMessage(user.ID, chat.ID))
		lastMessage, _ := conn.SetMessage(mocks.MakeMessage(mocks.ADMIN_ID, chat.ID))

		messages, firstMessageID, lastMessageID, httpError := getMessages(user.ID, chat.ID, 0, 0, 0, conn)
		assert.IsNil(t, httpError)
		assert.HasLength(t, messages, 3)
		assert.Equals(t, firstMessageID, firstMessage.ID)
		assert.Equals(t, lastMessageID, lastMessage.ID)
		// check ordering
		assert.Equals(t, messages[0].ID, firstMessage.ID)
		assert.Equals(t, messages[1].ID, middleMessage.ID)
		assert.Equals(t, messages[2].ID, lastMessage.ID)
	})
}

func TestSendMessageHTML(t *testing.T) {
	// TODO: perhaps test event sending
	t.Run("Normal", func(t *testing.T) {
		body := fmt.Sprintf(`{"content": "%s"}`, "This is a sample message!")
		w, r, conn := resolverutils.CommonSetup(body)
		user, _ := conn.SetUser(mocks.MakeUser())
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), user.ID, mocks.ADMIN_ID)
		params := map[string]string{resolverutils.CHAT_ID_KEY: fmt.Sprint(chat.ID)}
		r = resolverutils.SetContext(t, r, mocks.Admin, params)

		SendMessageHTML(w, r, conn)
		assert.Equals(t, w.Status, http.StatusNoContent)
		assert.Equals(t, string(w.Body), "")
	})

	t.Run("EmptyMessage", func(t *testing.T) {
		w, r, conn := resolverutils.CommonSetup(`{"content": ""}`)
		user, _ := conn.SetUser(mocks.MakeUser())
		chat, _ := conn.SetChat(mocks.MakePrivateChat(), user.ID, mocks.ADMIN_ID)
		params := map[string]string{resolverutils.CHAT_ID_KEY: fmt.Sprint(chat.ID)}
		r = resolverutils.SetContext(t, r, mocks.Admin, params)

		SendMessageHTML(w, r, conn)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Equals(t, string(w.Body), "cannot send an empty message")
	})
}

func TestUserSearch(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		body := fmt.Sprintf(`{"query": "%s"}`, "partial_username")
		w, r, conn := resolverutils.CommonSetup(body)
		r = resolverutils.SetContext(t, r, mocks.Admin, nil)

		UserSearch(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(t, string(w.Body), "<b>the_admin</b> Administrator")
	})
}

func TestCreatePrivateChatHTML(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		body := fmt.Sprintf(`{"userID": %d}`, mocks.ADMIN_ID)
		w, r, conn := resolverutils.CommonSetup(body)
		user, _ := conn.SetUser(mocks.MakeUser())
		r = resolverutils.SetContext(t, r, user, nil)

		CreatePrivateChatHTML(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Contains(
			t,
			string(w.Body),
			`<span class="heading-1">Administrator</span>`,
			`<div id=chat-list hx-swap-oob="innerHTML">`,
			`[private] <b>Administrator</b>`,
		)
	})
}

func TestRenameUser(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		xName := "Bukayo Saka"
		body := fmt.Sprintf(`{"newName": "%s"}`, xName)
		w, r, conn := resolverutils.CommonSetup(body)
		user, _ := conn.SetUser(mocks.MakeUser())
		r = resolverutils.SetContext(t, r, user, nil)

		RenameUser(w, r, conn)
		// formatting is weird to match spacing in original
		message := fmt.Sprintf(
			`<span class="info">
			Your display name has been changed to 
			<span class="accent">%s</span>.
			Your username is unchanged.
		</span>`,
			xName,
		)
		assert.Equals(t, w.Status, http.StatusOK)
		assert.Equals(t, string(w.Body), message)
	})
}
