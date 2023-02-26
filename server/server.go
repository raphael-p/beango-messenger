package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/raphael-p/beango/resolvers"
)

func Start() {
	router := newRouter()
	router.PUT("/session", resolvers.CreateSession)
	router.GET("/users", resolvers.GetUsers)
	router.POST("/user", resolvers.CreateUser)
	router.GET("/chats", resolvers.GetChats)
	router.PUT("/chat", resolvers.CreateChat)
	router.GET("/messages", resolvers.GetMessages)
	router.POST("/message", resolvers.SendMessage)

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}
	fmt.Println("🐱‍💻 BeanGo server started on", l.Addr().String())
	if err := http.Serve(l, router); err != nil {
		fmt.Printf("server closed: %s\n", err)
	}
	os.Exit(1)
}
