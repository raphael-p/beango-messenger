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
	defer utils.Logger.Close()

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
		utils.Logger.Error(fmt.Sprintf("failed to start server: %s\n", err))
	}
	utils.Logger.Info(fmt.Sprintf("🐱‍💻 BeanGo server started on %s\n", l.Addr().String()))
	if err := http.Serve(l, router); err != nil {
		utils.Logger.Error(fmt.Sprintf("server closed: %s\n", err))
	}
	os.Exit(1)
}
