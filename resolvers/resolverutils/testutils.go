package resolverutils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/response"
)

// Creates mock resolver arguments
func CommonSetup(body string) (*response.Writer, *http.Request, database.Connection) {
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := response.NewWriter(httptest.NewRecorder())
	return w, req, mocks.MakeMockConnection()
}

// Adds a user and/or parameters to a new request context
func SetContext(
	t *testing.T,
	req *http.Request,
	user *database.User,
	params map[string]string,
) *http.Request {
	var err error = nil
	if user != nil {
		req, err = context.SetUser(req, user)
		assert.IsNil(t, err)
	}
	for key, value := range params {
		req, err = context.SetParam(req, key, value)
		assert.IsNil(t, err)
	}
	return req
}

func AssertHTTPError(t *testing.T, err *HTTPError, expectedStatus int, expectedMessage string) {
	if err == nil {
		t.Error("expected HTTPError, got nil")
		return
	}
	if err.Status != expectedStatus {
		t.Errorf("expected HTTPError status \"%d\", got \"%d\"", expectedStatus, err.Status)
	}
	if err.Message != expectedMessage {
		t.Errorf("expected HTTPError message \"%s\", got \"%s\"", expectedMessage, err.Message)
	}
}
