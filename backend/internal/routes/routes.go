package routes

import (
	"github.com/ashmit-singh-gogia/c-hat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine, userHandler *handlers.UserHandler, chatHandler *handlers.ChatHandler) {

	api := router.Group("/api")
	users := api.Group("/users") // creates a user subgroup inside the api group
	users.POST("/", userHandler.RegisterUser)

	chats := api.Group("/chats")
	chats.POST("/direct", chatHandler.CreateDirectChat)
}
