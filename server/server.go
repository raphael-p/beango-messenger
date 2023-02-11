package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/raphael-p/beango/resolvers"
)

type loggingResponseWriter struct {
	status int
	body   string
	time   int64
	http.ResponseWriter
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	w.body = string(body)
	return w.ResponseWriter.Write(body)
}

func (w *loggingResponseWriter) String() string {
	out := fmt.Sprintf("status %d (took %dms)", w.status, w.time)
	if w.body != "" {
		out = fmt.Sprintf("%s\n\tresponse: %s", out, w.body)
	}
	return out
}

func logger(h http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	requestString := fmt.Sprint(r.Method, " ", r.URL)
	fmt.Println("received ", requestString)
	loggingResponseWriter := &loggingResponseWriter{ResponseWriter: w}
	start := time.Now()
	h(loggingResponseWriter, r)
	loggingResponseWriter.time = time.Since(start).Milliseconds()
	fmt.Printf(
		"%s resolved with %s\n",
		requestString,
		loggingResponseWriter,
	)
}

type route struct {
	method  string
	pattern *regexp.Regexp
	handler http.HandlerFunc
}

var routes []route

type ctxKey struct{}

func newRoute(method, pattern string, handler http.HandlerFunc) {
	route := route{method, regexp.MustCompile("^" + pattern + "$"), handler}
	routes = append(routes, route)
}

func router(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.pattern.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			logger(route.handler, w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func Start() {
	newRoute(http.MethodGet, "/users", resolvers.GetUsers)
	newRoute(http.MethodPost, "/user", resolvers.CreateUser)
	// "/chats", resolvers.GetChats
	// "/chat", resolvers.CreateChat
	// "/messages", resolvers.GetMessages
	// "/message", resolvers.SendMessage

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}
	fmt.Println("🐱‍💻 BeanGo server started on", l.Addr().String())
	if err := http.Serve(l, http.HandlerFunc(router)); err != nil {
		fmt.Printf("server closed: %s\n", err)
	}
	os.Exit(1)
}
