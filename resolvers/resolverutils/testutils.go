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

func CommonSetup(body string) (*response.Writer, *http.Request, database.Connection) {
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := response.NewWriter(httptest.NewRecorder())
	return w, req, mocks.MakeMockConnection()
}

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
