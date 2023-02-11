package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
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

func logger(h func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestString := fmt.Sprint(r.Method, " ", r.URL)
		fmt.Println("received ", requestString)
		loggingResponseWriter := &loggingResponseWriter{ResponseWriter: w}
		start := time.Now()
		http.HandlerFunc(h).ServeHTTP(loggingResponseWriter, r)
		loggingResponseWriter.time = time.Since(start).Milliseconds()
		fmt.Printf(
			"%s resolved with %s\n",
			requestString,
			loggingResponseWriter,
		)
	})

}

func Start() {
	http.Handle("/users", logger(resolvers.GetUsers))
	http.Handle("/user", logger(resolvers.CreateUser))
	// "/chats", resolvers.GetChats
	// "/chat", resolvers.CreateChat
	// "/messages", resolvers.GetMessages
	// "/message", resolvers.SendMessage

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}
	fmt.Println("🐱‍💻 BeanGo server started on", l.Addr().String())
	if err := http.Serve(l, nil); err != nil {
		fmt.Println("server closed")
	}
	os.Exit(1)
}
