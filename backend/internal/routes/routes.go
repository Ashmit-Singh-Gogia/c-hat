package routes

import (
	"os"

	"github.com/ashmit-singh-gogia/c-hat/internal/handlers"
	"github.com/ashmit-singh-gogia/c-hat/internal/middleware"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine, userHandler *handlers.UserHandler, chatHandler *handlers.ChatHandler, messageHandler *handlers.MessageHandler, authHandler *handlers.AuthHandler) {

	// CORS Middleware to allow React frontend to communicate with Gin backend
	router.Use(func(c *gin.Context) {
		// Change this if your Vite server runs on a different port
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := router.Group("/api")
	auth := api.Group("/auth")
	oauth := auth.Group("/oauth")
	// OAuth routes
	oauth.GET("/:provider", authHandler.HandleLogin)
	oauth.GET("/:provider/callback", authHandler.HandleCallback)
	// Local auth routes
	localAuth := auth.Group("/local")
	localAuth.POST("/register", authHandler.HandleLocalRegister)
	localAuth.POST("/login", authHandler.HandleLocalLogin)
	localAuth.POST("/resend-verification", authHandler.HandleResendVerification)
	localAuth.GET("/verify-email", authHandler.HandleEmailVerification)
	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		protected.GET("/auth/logout", authHandler.HandleLogout)
		// Fetch current logged-in user details
		protected.GET("/users/me", userHandler.GetMe)

		chats := protected.Group("/chats")
		chats.GET("/", chatHandler.GetChats) // Fetch all chats for sidebar
		chats.POST("/direct", chatHandler.CreateDirectChat)
		chats.GET("/:id/messages", messageHandler.GetMessages)

		messages := protected.Group("/messages")
		messages.POST("/", messageHandler.SendMessage)
	}
}
