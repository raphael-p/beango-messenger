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

func TestRefreshChat(t *testing.T) {
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

		RefreshChat(w, r, conn)
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

		RefreshChat(w, r, conn)
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
