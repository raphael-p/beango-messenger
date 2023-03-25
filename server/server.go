package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/utils"
)

func Start() {
	// Set up logger
	logger, err := utils.NewLogger("logs", "server.log")
	if err != nil {
		fmt.Printf("FATAL ERROR: %s\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Set up router
	router := newRouter()
	router.POST("/session", resolvers.CreateSession).noAuth()
	router.POST("/user", resolvers.CreateUser).noAuth()
	router.GET("/user/:username", resolvers.GetUserByName)
	router.GET("/chats", resolvers.GetChats)
	router.POST("/chat", resolvers.CreateChat)
	router.GET("/messages/:chatid", resolvers.GetChatMessages)
	router.POST("/message/:chatid", resolvers.SendMessage)

	// Run server
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
