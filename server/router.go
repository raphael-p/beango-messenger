package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/utils"
)

type handlerFunc func(*utils.ResponseWriter, *http.Request)

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
	paramKeys := []string{}
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
			req = buildContext(req, route.paramKeys, matches[1:])
			route.handler(utils.NewResponseWriter(w), req)
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

// Server's context keys, used to avoid clashes
type ContextParameters string
type ContextUser string

// Returns a shallow-copy of the request with an updated context,
// including path parameters
func buildContext(req *http.Request, paramKeys, paramValues []string) *http.Request {
	ctx := req.Context()
	for i := 0; i < len(paramKeys); i++ {
		ctx = context.WithValue(ctx, ContextParameters(paramKeys[i]), paramValues[i])
	}
	return req.WithContext(ctx)
}

// A wrapper around a route's handler for request middleware
func (r *route) handler(w *utils.ResponseWriter, req *http.Request) {
	// Log request
	requestString := fmt.Sprint(req.Method, " ", req.URL)
	fmt.Println("received ", requestString)

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
	fmt.Printf("%s resolved with %s\n", requestString, w)
}

func authentication(w *utils.ResponseWriter, req *http.Request) (*http.Request, error) {
	userId, err := getUserIdFromCookie(w, req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}
	user, err := database.GetUser(userId)
	if err != nil {
		w.StringResponse(http.StatusInternalServerError, "session is valid, but user no longer exists")
		return nil, err
	}
	contextWithUser := context.WithValue(req.Context(), utils.ContextUser("user"), user)
	return req.WithContext(contextWithUser), nil
}

func getUserIdFromCookie(w *utils.ResponseWriter, req *http.Request) (string, error) {
	cookieName := utils.AUTH_COOKIE
	sessionId, err := utils.GetCookieValue(cookieName, req)
	if err != nil {
		return "", err
	}
	session, ok := database.CheckSession(sessionId)
	if !ok {
		utils.InvalidateCookie(cookieName, w)
		return "", errors.New("cookie or session is invalid")
	}
	return session.UserId, nil
}
