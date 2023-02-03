package server

import (
	"github.com/gin-gonic/gin"
	"github.com/raphael-p/beango/resolvers"
)

func Start() {
	router := gin.Default()
	router.POST("/message", resolvers.SendMessage)
	router.PUT("/chat", resolvers.CreateChat)
	router.POST("/user", resolvers.CreateUser)
	router.GET("/messages", resolvers.GetMessages)
	router.GET("/chats", resolvers.GetChats)
	router.GET("/users", resolvers.GetUsers)
	router.Run("localhost:8080")
}
