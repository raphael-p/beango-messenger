package resolvers

import (
	"net/http"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/response"
	"github.com/raphael-p/beango/utils/validate"
)

func TestValidateCreateUserInput(t *testing.T) {
	makeInput := func(username, displayName, password string) createUserInput {
		return createUserInput{
			username,
			validate.JSONField[string]{Value: displayName},
			password,
		}
	}

	validateInput := func(t *testing.T, input createUserInput, xUsername, xDisplayName, xPassword string) {
		assert.Equals(t, input.Username, xUsername)
		assert.Equals(t, input.DisplayName.Value, xDisplayName)
		assert.Equals(t, input.Password, xPassword)

	}

	t.Run("Normal", func(t *testing.T) {
		xUsername := "gouser123"
		xDisplayName := "James Jameson"
		xPassword := "secretpass 123"
		input := makeInput(xUsername, xDisplayName, xPassword)

		err := validateCreateUserInput(&input)
		assert.IsNil(t, err)
		validateInput(t, input, xUsername, xDisplayName, xPassword)
	})

	t.Run("TrimsSpaces", func(t *testing.T) {
		xUsername := "gouser123"
		xDisplayName := "James Jameson"
		xPassword := "secretpass 123"
		input := makeInput(" \t"+xUsername+"\r", " "+xDisplayName+"\n", xPassword)

		err := validateCreateUserInput(&input)
		assert.IsNil(t, err)
		validateInput(t, input, xUsername, xDisplayName, xPassword)
	})

	t.Run("UsernameMaxLength", func(t *testing.T) {
		username := "abcdefghijklmnopqrstuvwxyz1"
		input := makeInput(username, "", "")

		err := validateCreateUserInput(&input)
		xMessage := "username must be shorter than 26 characters"
		resolverutils.AssertHTTPError(t, err, http.StatusBadRequest, xMessage)
	})

	t.Run("UsernameHasSpace", func(t *testing.T) {
		username := "iam auser"
		input := makeInput(username, "", "")

		err := validateCreateUserInput(&input)
		xMessage := "username may only contain alphanumeric characters and '_.'"
		resolverutils.AssertHTTPError(t, err, http.StatusBadRequest, xMessage)
	})

	t.Run("UsernameHasSpecialCharacter", func(t *testing.T) {
		username := "iam%auser"
		input := makeInput(username, "", "")

		err := validateCreateUserInput(&input)
		xMessage := "username may only contain alphanumeric characters and '_.'"
		resolverutils.AssertHTTPError(t, err, http.StatusBadRequest, xMessage)
	})

	t.Run("DisplayNameMaxLength", func(t *testing.T) {
		displayName := "abcdefghijklmnopqrstuvwxyz"
		input := makeInput("", displayName, "")

		err := validateCreateUserInput(&input)
		xMessage := "display name must be shorter than 26 characters"
		resolverutils.AssertHTTPError(t, err, http.StatusBadRequest, xMessage)
	})

	t.Run("DisplayNameHasWhitespace", func(t *testing.T) {
		displayName := "abcdefghij\rklmnopqrs"
		input := makeInput("", displayName, "")

		err := validateCreateUserInput(&input)
		xMessage := "display name may not contain any tabs or newlines"
		resolverutils.AssertHTTPError(t, err, http.StatusBadRequest, xMessage)
	})
}

func TestCreateUserDatabase(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		username := "xXbeanXx"
		display := "Bean The Cat"
		conn := mocks.MakeMockConnection()

		output, httpError := createUserDatabase(username, display, "abc123", conn)
		assert.IsNil(t, httpError)
		assert.Equals(t, output.Username, username)
		assert.Equals(t, output.DisplayName, display)

		user, err := conn.GetUser(output.ID)
		assert.IsNil(t, err)
		assert.IsNotNil(t, user)
		assert.HasLength(t, user.Key, 60) // typical bcrypt hash length
	})

	t.Run("UsernameTaken", func(t *testing.T) {
		conn := mocks.MakeMockConnection()

		output, httpError := createUserDatabase(mocks.ADMIN_USERNAME, "Bean", "abc123", conn)
		assert.IsNil(t, output)
		assert.Equals(t, httpError.Status, http.StatusConflict)
		assert.Equals(t, httpError.Message, "username is taken")
	})

	t.Run("PasswordNotHashable", func(t *testing.T) {
		password := "This is string is longer than 72 bytes. " +
			"bcrypt will not like this string."
		conn := mocks.MakeMockConnection()

		output, httpError := createUserDatabase("xXBeanXx", "Bean", password, conn)
		assert.IsNil(t, output)
		assert.Equals(t, httpError.Status, http.StatusBadRequest)
		assert.Equals(t, httpError.Message, "bcrypt: password length exceeds 72 bytes")
	})

	t.Run("NoDisplayName", func(t *testing.T) {
		username := "xXbeanXx"
		conn := mocks.MakeMockConnection()

		output, httpError := createUserDatabase(username, "", "abc123", conn)
		assert.IsNil(t, httpError)
		assert.Equals(t, output.DisplayName, username)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		body := `{"Username": "xXbeanXx", "displayName": "Bean The Cat", "password":"abc123"}`
		w, r, conn := resolverutils.CommonSetup(body)

		CreateUser(w, r, conn)
		assert.Equals(t, w.Status, http.StatusCreated)
		assert.IsValidJSON(t, string(w.Body), &userOutput{})
	})
}

func TestGetUserByName(t *testing.T) {
	setup := func(key, value string) (
		*response.Writer,
		*http.Request,
		database.Connection,
	) {
		w, r, conn := resolverutils.CommonSetup("")
		if key == "" {
			key = "username"
		}
		if value == "" {
			value = mocks.Admin.Username
		}
		params := map[string]string{key: value}
		r = resolverutils.SetContext(t, r, mocks.MakeUser(), params)
		return w, r, conn
	}

	t.Run("Normal", func(t *testing.T) {
		w, r, conn := setup("", "")

		GetUserByName(w, r, conn)
		assert.Equals(t, w.Status, http.StatusOK)
		xOutput := stripUserFields(*mocks.Admin)[0]
		var output userOutput
		assert.IsValidJSON(t, string(w.Body), &output)
		assert.Equals(t, output, xOutput)
	})

	t.Run("UsernameParamNotSet", func(t *testing.T) {
		w, r, conn := setup("not-username", "")

		GetUserByName(w, r, conn)
		assert.Equals(t, w.Status, http.StatusInternalServerError)
		assert.Equals(t, string(w.Body), "failed to fetch path parameter: username")
	})

	t.Run("NoMatchingUsername", func(t *testing.T) {
		w, r, conn := setup("", "xXbeanXx")

		GetUserByName(w, r, conn)
		assert.Equals(t, w.Status, http.StatusNotFound)
		assert.Equals(t, string(w.Body), "user not found")
	})
}
