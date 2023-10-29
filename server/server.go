package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/raphael-p/beango/config"
	"github.com/raphael-p/beango/database"
	"github.com/raphael-p/beango/resolvers"
	"github.com/raphael-p/beango/resolvers/resolverutils"
	"github.com/raphael-p/beango/server/routing"
	"github.com/raphael-p/beango/utils/logger"
	"github.com/raphael-p/beango/utils/path"
	"github.com/raphael-p/beango/utils/response"
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

	path, ok := path.RelativeJoin("../client/resources")
	if !ok {
		panic("failed to get path at runtime")
	}

	// aliases for readability
	chatID := resolverutils.CHAT_ID_KEY
	username := resolverutils.USERNAME_KEY

	// frontend endpoints
	router.GET("/login", resolvers.Login)
	router.POST("/login/:action", resolvers.SubmitLogin)
	router.GET("/home", resolvers.Home, routing.AuthRedirect)
	router.GET("/home/chat/:"+chatID, resolvers.OpenChat, routing.AuthRedirect)
	router.GET("/home/chat/:"+chatID+"/scrollUp", resolvers.ScrollUp, routing.AuthRedirect)
	router.GET("/home/chat/:"+chatID+"/refresh", resolvers.RefreshChat, routing.AuthRedirect)
	router.GET("/home/newChat", resolvers.OpenChatCreator, routing.AuthRedirect)
	router.POST("/home/newChat/search", resolvers.UserSearch, routing.AuthRedirect)
	router.POST("/home/newChat/create", resolvers.CreatePrivateChatHTML, routing.AuthRedirect)
	router.POST("/home/chat/:"+chatID+"/sendMessage", resolvers.SendMessageHTML, routing.AuthRedirect)
	router.GET("/favicon.ico", func(w *response.Writer, r *http.Request, conn database.Connection) {
		http.FileServer(http.Dir(path)).ServeHTTP(w, r)
	})
	router.GET("/resources/.*", func(w *response.Writer, r *http.Request, conn database.Connection) {
		http.StripPrefix("/resources/", http.FileServer(http.Dir(path))).ServeHTTP(w, r)
	})

	// backend endpoints
	router.POST("/session", resolvers.CreateSession)
	router.POST("/user", resolvers.CreateUser)
	router.GET("/user/:"+username, resolvers.GetUserByName, routing.Auth)
	router.GET("/chats", resolvers.GetChats, routing.Auth)
	router.POST("/chat", resolvers.CreatePrivateChat, routing.Auth)
	router.GET("/chat/:"+chatID+"/messages", resolvers.GetChatMessages, routing.Auth)
	router.POST("/chat/:"+chatID+"/message", resolvers.SendMessage, routing.Auth)

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
