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
	"github.com/raphael-p/beango/utils/validate"
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
	pathParamMatcher := regexp.MustCompile(":([a-zA-Z]+)")
	matches := pathParamMatcher.FindAllStringSubmatch(endpoint, -1)
	paramKeys := []string{}
	pattern := endpoint
	if len(matches) > 0 {
		// replace path parameter definition with regex pattern to capture any string
		pattern = pathParamMatcher.ReplaceAllLiteralString(endpoint, "([^/]+)")
		// store the names of path parameters, to later be used as context keys
		for i := 0; i < len(matches); i++ {
			paramKeys = append(paramKeys, matches[i][1])
		}
	}
	if !validate.UniqueList(paramKeys) {
		logger.Fatal(fmt.Sprint("duplicate path parameters in route: ", endpoint))
	}

	route := &route{
		method,
		regexp.MustCompile("^" + pattern + "$"),
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
				var err error
				req, err = context.SetParam(req, key, values[idx])
				if err != nil {
					logger.Error(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
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
		reqWithUser, ok := authentication(w, req)
		if !ok {
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

func authentication(w *response.Writer, req *http.Request) (*http.Request, bool) {
	userID, err := getUserIDFromCookie(w, req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, false
	}
	user, err := database.GetUser(userID)
	if err != nil {
		w.WriteString(http.StatusNotFound, "user not found during authentication")
		return nil, false
	}
	req, err = context.SetUser(req, user)
	if err != nil {
		logger.Error(err.Error())
		w.WriteString(http.StatusInternalServerError, err.Error())
		return nil, false
	}
	return req, true
}

func getUserIDFromCookie(w *response.Writer, req *http.Request) (string, error) {
	cookieName := cookies.SESSION
	sessionID, err := cookies.Get(req, cookieName)
	if err != nil {
		return "", err
	}
	session, ok := database.CheckSession(sessionID)
	if !ok {
		err := cookies.Invalidate(w, cookieName)
		if err != nil {
			logger.Error(err.Error())
		}
		return "", errors.New("cookie or session is invalid") // TODO: replace some fmt.Errorf with errors.New
	}
	return session.UserID, nil
}
