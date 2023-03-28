package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/utils"
)

// Handles failures on startup before the main logger can be created
func startupFailer(message string) {
	reset := "\033[0m"
	red := "\033[31;1m"
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger.Fatalf("%s[%s]%s startup failed: %s", red, "FATAL_ERROR", reset, message)
}

func Start() {
	utils.CreateConfig(startupFailer)
	utils.CreateLogger(startupFailer)
	defer utils.Logger.Close()

	router := newRouter()
	router.POST("/session", resolvers.CreateSession).noAuth()
	router.POST("/user", resolvers.CreateUser).noAuth()
	router.GET("/user/:username", resolvers.GetUserByName)
	router.GET("/chats", resolvers.GetChats)
	router.POST("/chat", resolvers.CreateChat)
	router.GET("/messages/:chatid", resolvers.GetChatMessages)
	router.POST("/message/:chatid", resolvers.SendMessage)

	l, err := net.Listen("tcp", fmt.Sprint(":", utils.Config.Server.Port))
	if err != nil {
		utils.Logger.Error(fmt.Sprint("failed to start server: ", err))
	}
	utils.Logger.Info(fmt.Sprintf("🐱‍💻 BeanGo server started on %s", l.Addr().String()))
	if err := http.Serve(l, router); err != nil {
		utils.Logger.Error(fmt.Sprintf("server closed: %s", err))
	}
	os.Exit(1)
}
