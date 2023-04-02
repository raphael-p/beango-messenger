package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/utils/logger"
)

func startupFailer(message string) {
	logger.Fatal(fmt.Sprint("startup failed: ", message))
}

func Start() {
	config.CreateConfig(startupFailer)
	logger.OpenLogFile(startupFailer)
	defer logger.CloseLogFile()

	router := newRouter()
	router.POST("/session", resolvers.CreateSession).noAuth()
	router.POST("/user", resolvers.CreateUser).noAuth()
	router.GET("/user/:username", resolvers.GetUserByName)
	router.GET("/chats", resolvers.GetChats)
	router.POST("/chat", resolvers.CreateChat)
	router.GET("/messages/:chatid", resolvers.GetChatMessages)
	router.POST("/message/:chatid", resolvers.SendMessage)
	l, err := net.Listen("tcp", fmt.Sprint(":", config.Values.Server.Port))
	if err != nil {
		logger.Error(fmt.Sprint("failed to start server: ", err))
	}
	logger.Info(fmt.Sprintf("🐱‍💻 BeanGo server started on %s", l.Addr().String()))
	if err := http.Serve(l, router); err != nil {
		logger.Error(fmt.Sprintf("server closed: %s", err))
	}
	os.Exit(1)
}
