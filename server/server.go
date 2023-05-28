package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/server/routing"
	"github.com/raphael-p/beango/utils/logger"
)

func setup() (router *routing.Router, ok bool) {
	ok = true
	defer func() {
		if r := recover(); r != nil {
			ok = false
			logger.Error(fmt.Sprint("setup failed: ", r))
		}
	}()

	config.CreateConfig()
	logger.Init()

	router = routing.NewRouter()
	router.POST("/session", resolvers.CreateSession).NoAuth()
	router.POST("/user", resolvers.CreateUser).NoAuth()
	router.GET("/user/:"+resolvers.USERNAME_KEY, resolvers.GetUserByName)
	router.GET("/chats", resolvers.GetChats)
	router.POST("/chat", resolvers.CreateChat)
	router.GET("/messages/:"+resolvers.CHAT_ID_KEY, resolvers.GetChatMessages)
	router.POST("/message/:"+resolvers.CHAT_ID_KEY, resolvers.SendMessage)
	return router, ok
}

func Start() {
	defer os.Exit(1)
	defer logger.Close()
	router, ok := setup()
	if !ok {
		return
	}

	l, err := net.Listen("tcp", fmt.Sprint(":", config.Values.Server.Port))
	if err != nil {
		logger.Error(fmt.Sprint("failed to start server: ", err))
		return
	}
	logger.Info(fmt.Sprintf("🐱‍💻 BeanGo server started on %s", l.Addr().String()))
	if err := http.Serve(l, router); err != nil {
		logger.Error(fmt.Sprint("server closed: ", err))
	}
}
