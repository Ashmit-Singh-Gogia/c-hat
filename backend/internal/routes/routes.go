package routes

import (
	"os"

	"github.com/ashmit-singh-gogia/c-hat/internal/handlers"
	"github.com/ashmit-singh-gogia/c-hat/internal/middleware"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine, userHandler *handlers.UserHandler, chatHandler *handlers.ChatHandler, messageHandler *handlers.MessageHandler, authHandler *handlers.AuthHandler) {

	api := router.Group("/api")
	// Google OAuth routes
	api.GET("/auth/google", authHandler.GoogleLogin)
	api.GET("/auth/google/callback", authHandler.GoogleCallback)

	// User registration route
	users := api.Group("/users") // creates a user subgroup inside the api group
	users.POST("/", userHandler.RegisterUser)

	// main.go or routes.go
	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		chats := protected.Group("/chats")
		chats.POST("/direct", chatHandler.CreateDirectChat)
		chats.GET("/:id/messages", messageHandler.GetMessages)

		messages := protected.Group("/messages")
		messages.POST("/", messageHandler.SendMessage)
	}
}
