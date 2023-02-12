package server

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/raphael-p/beango/utils"
)

type handlerFunc func(*utils.ResponseWriter, *http.Request)

type route struct {
	method  string
	pattern *regexp.Regexp
	handler handlerFunc
}

type router struct {
	routes []route
}

func newRouter() *router {
	return &router{routes: []route{}}
}

func (r *router) newRoute(method, pattern string, handler handlerFunc) {
	route := route{method, regexp.MustCompile("^" + pattern + "$"), handler}
	r.routes = append(r.routes, route)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.serveHTTP(utils.NewResponseWriter(w), req)
}

func (r *router) serveHTTP(w *utils.ResponseWriter, req *http.Request) {
	var allow []string
	for _, route := range r.routes {
		matches := route.pattern.FindStringSubmatch(req.URL.Path)
		if len(matches) > 0 {
			if req.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(req.Context(), struct{}{}, matches[1:])
			logger(route.handler, w, req.WithContext(ctx))
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

func logger(h handlerFunc, w *utils.ResponseWriter, r *http.Request) {
	requestString := fmt.Sprint(r.Method, " ", r.URL)
	fmt.Println("received ", requestString)
	start := time.Now()
	h(w, r)
	w.Time = time.Since(start).Milliseconds()
	fmt.Printf("%s resolved with %s\n", requestString, w)
}
