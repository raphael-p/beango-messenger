package resolvers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
)

func TestCreateUser(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		username := "xXbeanXx"
		display := "Bean The Cat"
		body := fmt.Sprintf(
			`{"Username": "%s", "displayName": "%s", "password":"abc123"}`,
			username,
			display,
		)
		w, req := mockRequest(body)
		conn := mocks.MakeMockConnection()

		CreateUser(w, req, conn)
		assert.Equals(t, w.Status, http.StatusCreated)

		var output UserOutput
		assert.IsValidJSON(t, w.Body, &output)
		assert.Equals(t, output.Username, username)
		assert.Equals(t, output.DisplayName, display)

		user, err := conn.GetUser(output.ID)
		assert.IsNil(t, err)
		assert.IsNotNil(t, user)
		assert.HasLength(t, user.Key, 60) // typical bcrypt hash length
	})

	t.Run("UsernameTaken", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"Username": "%s", "displayName": "Bean", "password":"abc123"}`,
			mocks.ADMIN_USERNAME,
		)
		w, req := mockRequest(body)
		conn := mocks.MakeMockConnection()

		CreateUser(w, req, conn)
		assert.Equals(t, w.Status, http.StatusConflict)
		assert.Equals(t, w.Body, "username is taken")
	})

	t.Run("PasswordNotHashable", func(t *testing.T) {
		password := "This is string is longer than 72 bytes. bcrypt will not like this string."
		body := fmt.Sprintf(`{"Username": "xXbeanXx", "displayName": "Bean", "password":"%s"}`, password)
		w, req := mockRequest(body)
		conn := mocks.MakeMockConnection()

		CreateUser(w, req, conn)
		assert.Equals(t, w.Status, http.StatusBadRequest)
		assert.Equals(t, w.Body, "bcrypt: password length exceeds 72 bytes")
	})

	t.Run("NoDisplayName", func(t *testing.T) {
		username := "xXbeanXx"
		body := fmt.Sprintf(`{"Username": "%s", "password":"abc123"}`, username)
		w, req := mockRequest(body)
		conn := mocks.MakeMockConnection()

		CreateUser(w, req, conn)
		assert.Equals(t, w.Status, http.StatusCreated)
		var output UserOutput
		assert.IsValidJSON(t, w.Body, &output)
		assert.Equals(t, output.DisplayName, username)
	})
}

func TestGetUserByName(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		w, req := mockRequest("")
		conn := mocks.MakeMockConnection()
		params := map[string]string{"username": mocks.Admin.Username}
		req = setContext(t, req, mocks.MakeUser(), params)

		GetUserByName(w, req, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		xOutput := *stripFields(mocks.Admin)
		var output UserOutput
		assert.IsValidJSON(t, w.Body, &output)
		assert.Equals(t, output, xOutput)
	})

	t.Run("UsernameParamNotSet", func(t *testing.T) {
		w, req := mockRequest("")
		conn := mocks.MakeMockConnection()
		params := map[string]string{"not-username": mocks.Admin.Username}
		req = setContext(t, req, mocks.MakeUser(), params)

		GetUserByName(w, req, conn)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, w.Body, "failed to fetch path parameter: username")
	})

	t.Run("NoMatchingUsername", func(t *testing.T) {
		w, req := mockRequest("")
		conn := mocks.MakeMockConnection()
		params := map[string]string{"username": "xXbeanXx"}
		req = setContext(t, req, mocks.MakeUser(), params)

		GetUserByName(w, req, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, w.Body, "user not found")
	})
}
