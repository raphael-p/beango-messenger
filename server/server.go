package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/server/routing"
	"github.com/raphael-p/beango/utils/logger"
)

func setup() (conn *database.MongoConnection, router *routing.Router, ok bool) {
	ok = true
	defer func() {
		if r := recover(); r != nil {
			ok = false
			logger.Error(fmt.Sprint("failed setup: ", r))
		}
	}()

	config.CreateConfig()
	logger.Init()

	conn, err := database.GetConnection()
	if err != nil {
		panic("failed to open database connection: " + err.Error())
	}
	logger.Trace("opened database connection")
	database.Setup(conn)

	router = routing.NewRouter()
	router.POST("/session", resolvers.CreateSession).NoAuth()
	router.POST("/user", resolvers.CreateUser).NoAuth()
	router.GET("/user/:"+resolvers.USERNAME_KEY, resolvers.GetUserByName)
	router.GET("/chats", resolvers.GetChats)
	router.POST("/chat", resolvers.CreateChat)
	router.GET("/messages/:"+resolvers.CHAT_ID_KEY, resolvers.GetChatMessages)
	router.POST("/message/:"+resolvers.CHAT_ID_KEY, resolvers.SendMessage)
	return conn, router, ok
}

func teardown(conn *database.MongoConnection) {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			logger.Error("failed to close database connection: " + err.Error())
		} else {
			logger.Trace("closed database connection")
		}
	}
	logger.Close()
	os.Exit(1)
}

func Start() {
	conn, router, ok := setup()
	defer teardown(conn)
	if !ok {
		return
	}

	l, err := net.Listen("tcp", fmt.Sprint(":", config.Values.Server.Port))
	if err != nil {
		logger.Error(fmt.Sprint("failed to start server: ", err))
		return
	}
	logger.Info(fmt.Sprintf("🐱‍💻 started BeanGo server on %s", l.Addr().String()))
	if err := http.Serve(l, router); err != nil {
		logger.Error(fmt.Sprint("closed server: ", err))
	}
}
