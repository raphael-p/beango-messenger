package resolvers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/response"
)

func TestCreateUser(t *testing.T) {
	setup := func(name, display, pass string) (
		*response.Writer,
		*http.Request,
		database.Connection,
	) {
		var body string
		if display == "" {
			body = fmt.Sprintf(
				`{"Username": "%s", "password":"%s"}`,
				name,
				pass,
			)
		} else {
			body = fmt.Sprintf(
				`{"Username": "%s", "displayName": "%s", "password":"%s"}`,
				name,
				display,
				pass,
			)
		}
		return resolverutils.CommonSetup(body)
	}

	t.Run("Normal", func(t *testing.T) {
		username := "xXbeanXx"
		display := "Bean The Cat"
		w, req, conn := setup(username, display, "abc123")

		CreateUser(w, req, conn)
		assert.Equals(t, w.Status, http.StatusCreated)

		var output UserOutput
		assert.IsValidJSON(t, string(w.Body), &output)
		assert.Equals(t, output.Username, username)
		assert.Equals(t, output.DisplayName, display)

		user, err := conn.GetUser(output.ID)
		assert.IsNil(t, err)
		assert.IsNotNil(t, user)
		assert.HasLength(t, user.Key, 60) // typical bcrypt hash length
	})

	t.Run("UsernameTaken", func(t *testing.T) {
		w, req, conn := setup(mocks.ADMIN_USERNAME, "Bean", "abc123")

		CreateUser(w, req, conn)
		assert.Equals(t, w.Status, http.StatusConflict)
		assert.Equals(t, string(w.Body), "username is taken")
	})

	t.Run("PasswordNotHashable", func(t *testing.T) {
		password := "This is string is longer than 72 bytes. " +
			"bcrypt will not like this string."
		w, req, conn := setup("xXBeanXx", "Bean", password)

		CreateUser(w, req, conn)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Equals(t, string(w.Body), "bcrypt: password length exceeds 72 bytes")
	})

	t.Run("NoDisplayName", func(t *testing.T) {
		username := "xXbeanXx"
		w, req, conn := setup(username, "", "abc123")

		CreateUser(w, req, conn)
		assert.Equals(t, w.Status, http.StatusCreated)
		var output UserOutput
		assert.IsValidJSON(t, string(w.Body), &output)
		assert.Equals(t, output.DisplayName, username)
	})
}

func TestGetUserByName(t *testing.T) {
	setup := func(key, value string) (
		*response.Writer,
		*http.Request,
		database.Connection,
	) {
		w, req, conn := resolverutils.CommonSetup("")
		if key == "" {
			key = "username"
		}
		if value == "" {
			value = mocks.Admin.Username
		}
		params := map[string]string{key: value}
		req = resolverutils.SetContext(t, req, mocks.MakeUser(), params)
		return w, req, conn
	}

	t.Run("Normal", func(t *testing.T) {
		w, req, conn := setup("", "")

		GetUserByName(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		xOutput := *stripFields(mocks.Admin)
		var output UserOutput
		assert.IsValidJSON(t, string(w.Body), &output)
		assert.Equals(t, output, xOutput)
	})

	t.Run("UsernameParamNotSet", func(t *testing.T) {
		w, req, conn := setup("not-username", "")

		GetUserByName(w, req, conn)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, string(w.Body), "failed to fetch path parameter: username")
	})

	t.Run("NoMatchingUsername", func(t *testing.T) {
		w, req, conn := setup("", "xXbeanXx")

		GetUserByName(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, string(w.Body), "user not found")
	})
}
