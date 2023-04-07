package mocks

import (
	"net/http/httptest"

	"github.com/raphael-p/beango/utils/response"
)

func MakeResponseWriter() *response.Writer {
	return &response.Writer{
		ResponseWriter: httptest.NewRecorder(),
	}
}
