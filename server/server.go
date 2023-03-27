package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/utils"
)

func Start() {
	// Start up config
	config.Init()

	// Start up logger
	utils.StartLogger()
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
	port := 8081 // TODO: make config variable
	l, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		utils.Logger.Error(fmt.Sprint("failed to start server: ", err))
	}
	utils.Logger.Info(fmt.Sprintf("🐱‍💻 BeanGo server started on %s", l.Addr().String()))
	if err := http.Serve(l, router); err != nil {
		utils.Logger.Error(fmt.Sprintf("server closed: %s", err))
	}
	os.Exit(1)
}
