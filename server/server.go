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
	router.POST("/session", resolvers.CreateSession).noAuth()
	router.GET("/users", resolvers.GetUsers)
	router.POST("/user", resolvers.CreateUser).noAuth()
	router.GET("/chats", resolvers.GetChats)
	router.POST("/chat", resolvers.CreateChat)
	router.GET("/messages/:chatid", resolvers.GetChatMessages)
	router.POST("/message/:chatid", resolvers.SendMessage)

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
