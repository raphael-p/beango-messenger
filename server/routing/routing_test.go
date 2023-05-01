package routing

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/test/assert"
	"github.com/raphael-p/beango/test/mocks"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

func TestNewRouter(t *testing.T) {
	r := NewRouter()
	assert.IsNotNil(t, r)
	assert.HasLength(t, r.routes, 0)
}

func TestNoAuth(t *testing.T) {
	r := &route{authenticate: true}
	r.NoAuth()
	assert.Equals(t, r.authenticate, false)
}

func makeRoute(method, pattern string, handler handlerFunc, params []string, auth bool) *route {
	return &route{
		method,
		regexp.MustCompile(pattern),
		handler,
		params,
		auth,
	}
}

func ptrAddress(ptr any) string { return fmt.Sprintf("%p", ptr) }

func assertRoute(t *testing.T, route *route, xRoute *route) {
	assert.Equals(t, route.method, xRoute.method)
	assert.Equals(t, route.pattern.String(), xRoute.pattern.String())
	assert.Equals(t, ptrAddress(route.innerHandler), ptrAddress(xRoute.innerHandler))
	assert.DeepEquals(t, route.paramKeys, xRoute.paramKeys)
	assert.Equals(t, route.authenticate, xRoute.authenticate)
}

func TestAddRoute(t *testing.T) {
	handler := func(w *response.Writer, r *http.Request, conn database.Connection) {}

	t.Run("Normal", func(t *testing.T) {
		router := NewRouter()
		method := http.MethodGet

		route := router.addRoute(method, "/from/:place/to/:newPlace/move", handler)
		xPattern := "^/from/([^/]+)/to/([^/]+)/move$"
		xParams := []string{"place", "newPlace"}
		xRoute := makeRoute(method, xPattern, handler, xParams, true)
		assertRoute(t, route, xRoute)
		assert.HasLength(t, router.routes, 1)
		assertRoute(t, router.routes[0], xRoute)
	})

	t.Run("DuplicateParamKeys", func(t *testing.T) {
		pathDef := "/foo/:id/bar/:id/b"
		defer func() {
			reason, ok := recover().(string)
			assert.Equals(t, ok, true)
			xReason := fmt.Sprint("duplicate parameters in path definition: ", pathDef)
			assert.Equals(t, reason, xReason)
		}()

		NewRouter().addRoute(http.MethodGet, pathDef, handler)
	})

	t.Run("NormalNoParams", func(t *testing.T) {

		route := NewRouter().addRoute(http.MethodGet, "/path/with/no/params", handler)
		assert.Equals(t, route.pattern.String(), "^/path/with/no/params$")
	})

	t.Run("InvalidParamNames", func(t *testing.T) {
		route := NewRouter().addRoute(http.MethodGet, "/path/with/:foo.bar/params", handler)
		assert.Equals(t, route.pattern.String(), "^/path/with/([^/]+).bar/params$")
		assert.DeepEquals(t, route.paramKeys, []string{"foo"})
	})

	t.Run("PatternMatchesRequestPath", func(t *testing.T) {
		path := "/path/with/123/param"

		route := NewRouter().addRoute(http.MethodGet, "/path/with/:num/param", handler)
		matches := route.pattern.FindStringSubmatch(path)
		assert.DeepEquals(t, matches, []string{path, "123"})
	})

	t.Run("Wrappers", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			router := NewRouter()
			pathDef := "/just/some/path"
			xPattern := fmt.Sprint("^", pathDef, "$")
			testCases := []struct {
				method          string
				addRouteWrapper func(string, handlerFunc) *route
			}{
				{http.MethodGet, router.GET},
				{http.MethodPost, router.POST},
				{http.MethodPut, router.PUT},
				{http.MethodPatch, router.PATCH},
				{http.MethodDelete, router.DELETE},
			}

			for idx, testCase := range testCases {
				t.Run(testCase.method, func(t *testing.T) {
					route := testCase.addRouteWrapper(pathDef, handler)
					xRoute := makeRoute(testCase.method, xPattern, handler, []string{}, true)
					assertRoute(t, route, xRoute)
					assert.HasLength(t, router.routes, idx+1)
					assertRoute(t, router.routes[idx], xRoute)
				})
			}
		})
	})
}

func TestServeHTTP(t *testing.T) {
	method := http.MethodGet
	pattern := "^/user/([^/]+)/name/([^/]+)$"
	path := func(id, name string) string {
		return fmt.Sprintf("/user/%s/name/%s", id, name)
	}
	code := http.StatusAccepted
	params := []string{"id", "name"}
	xBody := func(id, name, body string) string {
		return fmt.Sprintf("user: %s, name: %s, body: %s", id, name, body)
	}
	handler := func(w *response.Writer, req *http.Request, conn database.Connection) {
		id, err := context.GetParam(req, params[0])
		assert.IsNil(t, err)
		name, err := context.GetParam(req, params[1])
		assert.IsNil(t, err)
		body, err := io.ReadAll(req.Body)
		assert.IsNil(t, err)
		w.WriteString(code, xBody(id, name, string(body)))
	}

	newRoute := makeRoute(method, pattern, handler, params, false)
	router := Router{[]*route{newRoute}}
	origRoutes := append([]*route{}, router.routes...)
	resetRoutes := func() { router.routes = append([]*route{}, origRoutes...) }

	t.Run("ValidRequest", func(t *testing.T) {
		id := "19"
		name := "patrick"
		reqBody := `{"key": "value"}`
		bodyBuf := bytes.NewBufferString(reqBody)
		req := httptest.NewRequest(method, path(id, name), bodyBuf)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equals(t, res.Code, code)
		assert.Equals(t, res.Body.String(), xBody(id, name, reqBody))
	})

	t.Run("PicksCorrectRoute", func(t *testing.T) {
		xBody := "correct route picked"
		correctHandler := func(w *response.Writer, req *http.Request, conn database.Connection) {
			w.WriteString(code, xBody)
		}
		router.routes = append(
			router.routes,
			makeRoute(method, "^/correct$", correctHandler, []string{}, false),
		)
		defer resetRoutes()
		req := httptest.NewRequest(method, "/correct", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equals(t, res.Code, code)
		assert.Equals(t, res.Body.String(), xBody)
	})

	t.Run("PicksCorrectMethod", func(t *testing.T) {
		xBody := "correct method picked"
		correctMethod := http.MethodPatch
		correctHandler := func(w *response.Writer, req *http.Request, conn database.Connection) {
			w.WriteString(code, xBody)
		}
		router.routes = append(
			router.routes,
			makeRoute(correctMethod, "^/$", correctHandler, []string{}, false),
		)
		defer resetRoutes()
		req := httptest.NewRequest(correctMethod, "/", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equals(t, res.Code, code)
		assert.Equals(t, res.Body.String(), xBody)
	})

	t.Run("RunsAuth", func(t *testing.T) {
		oldAuth := newRoute.authenticate
		defer func() { newRoute.authenticate = oldAuth }()
		newRoute.authenticate = true
		req := httptest.NewRequest(method, path("3", "bean"), nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equals(t, res.Code, http.StatusUnauthorized)
		assert.Equals(t, res.Body.String(), "")
	})

	t.Run("DuplicateParamKey", func(t *testing.T) {
		oldParams := append([]string{}, newRoute.paramKeys...)
		defer func() { newRoute.paramKeys = oldParams }()
		newRoute.paramKeys = []string{"dupe", "dupe"}
		req := httptest.NewRequest(method, path("3", "bean"), nil)
		res := httptest.NewRecorder()
		buf := logger.MockFileLogger(t)

		router.ServeHTTP(res, req)
		assert.Equals(t, res.Code, http.StatusInternalServerError)
		xErr := "path parameter dupe already set"
		assert.Equals(t, res.Body.String(), xErr)
		assert.Contains(t, buf.String(), fmt.Sprint("[ERROR] ", xErr))
	})

	t.Run("WrongParamKeyCount", func(t *testing.T) {
		oldParams := append([]string{}, newRoute.paramKeys...)
		defer func() { newRoute.paramKeys = oldParams }()
		newRoute.paramKeys = append(newRoute.paramKeys, "extra")
		xPath := path("3", "bean")
		req := httptest.NewRequest(method, xPath, nil)
		res := httptest.NewRecorder()
		buf := logger.MockFileLogger(t)

		router.ServeHTTP(res, req)
		assert.Equals(t, res.Code, http.StatusInternalServerError)
		xErr := "unexpected number of path parameters in request"
		xErrLog := fmt.Sprintf("[ERROR] %s (%s)", xErr, xPath)
		assert.Equals(t, res.Body.String(), xErr)
		assert.Contains(t, buf.String(), xErrLog)
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		patchRoute := *newRoute
		patchRoute.method = http.MethodPatch
		router.routes = append(router.routes, &patchRoute)
		defer resetRoutes()
		req := httptest.NewRequest(http.MethodPost, path("3", "bean"), nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equals(t, res.Code, http.StatusMethodNotAllowed)
		assert.Equals(t, res.Body.String(), "")
		assert.Equals(t, res.Header().Get("Allow"), "GET, PATCH")
	})

	t.Run("PathNotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equals(t, res.Code, http.StatusNotFound)
	})
}

func TestRouteHandler(t *testing.T) {
	method := http.MethodGet
	path := "/"
	status := http.StatusOK
	body := "success"
	handler := func(w *response.Writer, r *http.Request, conn database.Connection) {
		w.WriteString(status, body)
	}
	params := []string{}
	authRoute := makeRoute(method, path, handler, params, true)
	noAuthRoute := makeRoute(method, path, handler, params, false)
	conn := mocks.MakeMockConnection(t)

	assertSuccess := func(t *testing.T, w *response.Writer, buf *bytes.Buffer) {
		assert.Equals(t, w.Status, status)
		assert.Equals(t, w.Body, body)
		reqString := fmt.Sprint(method, " ", path)
		receivedMessage := fmt.Sprint("[INFO] received ", reqString)
		successMessage := fmt.Sprintf("[INFO] %s resolved with %s", reqString, w)
		assert.Contains(t, buf.String(), receivedMessage, successMessage)
	}

	t.Run("AuthSucceeds", func(t *testing.T) {
		req := httptest.NewRequest(method, path, nil)
		cookie := &http.Cookie{Name: string(cookies.SESSION), Value: mocks.AdminSesh.ID}
		req.AddCookie(cookie)
		w := response.NewWriter(httptest.NewRecorder())
		buf := logger.MockFileLogger(t)

		authRoute.handler(w, req, conn)
		assertSuccess(t, w, buf)
	})

	t.Run("AuthFails", func(t *testing.T) {
		req := httptest.NewRequest(method, path, nil)
		w := response.NewWriter(httptest.NewRecorder())
		buf := logger.MockFileLogger(t)

		authRoute.handler(w, req, conn)
		assert.Equals(t, w.Status, http.StatusUnauthorized)
		receivedMessage := fmt.Sprint("[INFO] received ", method, " ", path)
		logs := buf.String()
		assert.Contains(t, logs, receivedMessage)
		assert.NotContains(t, logs, "resolved")
	})

	t.Run("AuthSkipped", func(t *testing.T) {
		req := httptest.NewRequest(method, path, nil)
		w := response.NewWriter(httptest.NewRecorder())
		buf := logger.MockFileLogger(t)

		noAuthRoute.handler(w, req, conn)
		assertSuccess(t, w, buf)
	})
}
