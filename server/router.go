package server

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/cookies"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
)

type handlerFunc func(*response.Writer, *http.Request)

type route struct {
	method       string
	pattern      *regexp.Regexp
	innerHandler handlerFunc
	paramKeys    []string
	authenticate bool
}

func (r *route) noAuth() {
	r.authenticate = false
}

type router struct {
	routes []*route
}

func newRouter() *router {
	return &router{routes: []*route{}}
}

func (r *router) addRoute(method, endpoint string, handler handlerFunc) *route {
	// handle path parameters
	pathParamPattern := regexp.MustCompile(":([a-z]+)")
	matches := pathParamPattern.FindAllStringSubmatch(endpoint, -1)
	paramKeys := []string{} // TODO: prevent duplicate keys
	if len(matches) > 0 {
		// replace path parameter definition with regex pattern to capture any string
		endpoint = pathParamPattern.ReplaceAllLiteralString(endpoint, "([^/]+)")
		// store the names of path parameters, to later be used as context keys
		for i := 0; i < len(matches); i++ {
			paramKeys = append(paramKeys, matches[i][1])
		}
	}

	route := &route{
		method,
		regexp.MustCompile("^" + endpoint + "$"),
		handler,
		paramKeys,
		true,
	}
	r.routes = append(r.routes, route)
	return route
}

func (r *router) GET(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodGet, pattern, handler)
}

func (r *router) POST(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodPost, pattern, handler)
}

func (r *router) PUT(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodPut, pattern, handler)
}

func (r *router) PATCH(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodPatch, pattern, handler)
}

func (r *router) DELETE(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodDelete, pattern, handler)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var allow []string
	for _, route := range r.routes {
		matches := route.pattern.FindStringSubmatch(req.URL.Path)
		if len(matches) > 0 {
			if req.Method != route.method {
				allow = append(allow, route.method)
				continue
			}

			values := matches[1:]
			if len(values) != len(route.paramKeys) {
				message := "unexpected number of path parameters in request"
				logger.Error(fmt.Sprint(message, " (", req.URL.Path, ")"))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(message))
				return
			}
			for idx, key := range route.paramKeys {
				req = context.SetParam(req, key, values[idx])
			}

			route.handler(response.NewWriter(w), req)
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, req)
}

// A wrapper around a route's handler for request middleware
func (r *route) handler(w *response.Writer, req *http.Request) {
	// Log request
	requestString := fmt.Sprint(req.Method, " ", req.URL)
	logger.Info(fmt.Sprint("received ", requestString))

	// Authentication
	if r.authenticate {
		reqWithUser, err := authentication(w, req)
		if err != nil {
			return
		}
		req = reqWithUser
	}

	// Log response
	start := time.Now()
	r.innerHandler(w, req)
	w.Time = time.Since(start).Milliseconds()
	logger.Info(fmt.Sprintf("%s resolved with %s", requestString, w))
}

func authentication(w *response.Writer, req *http.Request) (*http.Request, error) {
	userId, err := getUserIdFromCookie(w, req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}
	user, err := database.GetUser(userId)
	if err != nil {
		w.WriteString(http.StatusNotFound, "user not found during authentication")
		return nil, err
	}
	return context.SetUser(req, user), nil
}

func getUserIdFromCookie(w *response.Writer, req *http.Request) (string, error) {
	cookieName := cookies.SESSION
	sessionId, err := cookies.Get(req, cookieName)
	if err != nil {
		return "", err
	}
	session, ok := database.CheckSession(sessionId)
	if !ok {
		cookies.Invalidate(w, cookieName)
		return "", errors.New("cookie or session is invalid")
	}
	return session.UserId, nil
}
