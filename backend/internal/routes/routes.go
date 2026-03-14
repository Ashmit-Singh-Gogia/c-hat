package routes

import (
	"github.com/ashmit-singh-gogia/c-hat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine, userHandler *handlers.UserHandler, chatHandler *handlers.ChatHandler, messageHandler *handlers.MessageHandler) {

	api := router.Group("/api")
	users := api.Group("/users") // creates a user subgroup inside the api group
	users.POST("/", userHandler.RegisterUser)

	chats := api.Group("/chats")
	chats.POST("/direct", chatHandler.CreateDirectChat)
	chats.GET("/:id/messages", messageHandler.GetMessages)

	messages := api.Group("/messages")
	messages.POST("/", messageHandler.SendMessage)
}
