package routing

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/server/authenticate"
	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/response"
	"github.com/raphael-p/beango/utils/validate"
)

type handlerFunc func(*response.Writer, *http.Request, database.Connection)

type route struct {
	method       string
	pattern      *regexp.Regexp
	innerHandler handlerFunc
	paramKeys    []string
	authenticate bool
}

func (r *route) NoAuth() {
	r.authenticate = false
}

type Router struct {
	routes []*route
}

func NewRouter() *Router {
	return &Router{routes: []*route{}}
}

func (r *Router) addRoute(method, pathDef string, handler handlerFunc) *route {
	// handle path parameters
	pathParamMatcher := regexp.MustCompile(":([a-zA-Z]+)")
	matches := pathParamMatcher.FindAllStringSubmatch(pathDef, -1)
	paramKeys := []string{}
	pattern := pathDef
	if len(matches) > 0 {
		// replace path parameter definition with regex pattern to capture any string
		pattern = pathParamMatcher.ReplaceAllLiteralString(pathDef, "([^/]+)")
		// store the names of path parameters, to later be used as context keys
		for i := 0; i < len(matches); i++ {
			paramKeys = append(paramKeys, matches[i][1])
		}
	}
	if !validate.UniqueList(paramKeys) {
		panic(fmt.Sprint("duplicate parameters in path definition: ", pathDef))
	}

	// check for duplicates: same method and regex pattern
	regex := regexp.MustCompile("^" + pattern + "$")
	for _, route := range r.routes {
		if route.method == method && route.pattern.String() == regex.String() {
			panic(fmt.Sprintf("route already exists: %s %s", method, pathDef))
		}
	}

	newRoute := &route{
		method,
		regex,
		handler,
		paramKeys,
		true,
	}
	r.routes = append(r.routes, newRoute)
	return newRoute
}

func (r *Router) GET(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodGet, pattern, handler)
}

func (r *Router) POST(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodPost, pattern, handler)
}

func (r *Router) PUT(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodPut, pattern, handler)
}

func (r *Router) PATCH(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodPatch, pattern, handler)
}

func (r *Router) DELETE(pattern string, handler handlerFunc) *route {
	return r.addRoute(http.MethodDelete, pattern, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
				errorResponse := "unexpected number of path parameters in request"
				logger.Error(fmt.Sprintf("%s (%s)", errorResponse, req.URL.Path))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorResponse))
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

			conn, err := database.GetConnection()
			if err != nil {
				message := "failed to get database connection"
				logger.Error(message + ": " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(message))
			}

			route.handler(response.NewWriter(w), req, conn)
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
func (r *route) handler(w *response.Writer, req *http.Request, conn database.Connection) {
	// Log request
	requestString := fmt.Sprint(req.Method, " ", req.URL)
	logger.Info(fmt.Sprint("received ", requestString))

	// Authentication
	ok := false
	if r.authenticate {
		if req, ok = authenticate.FromCookie(w, req, conn); !ok {
			w.Commit()
			return
		}
	}

	// Log response
	start := time.Now()
	r.innerHandler(w, req, conn)
	w.Commit()
	w.Time = time.Since(start).Milliseconds()
	logger.Info(fmt.Sprintf("%s resolved with %s", requestString, w))
}
