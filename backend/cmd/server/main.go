package main

import (
	"fmt"

	"github.com/ashmit-singh-gogia/c-hat/internal/config"
	"github.com/ashmit-singh-gogia/c-hat/internal/database"
	"github.com/ashmit-singh-gogia/c-hat/internal/handlers"
	"github.com/ashmit-singh-gogia/c-hat/internal/repositories"
	"github.com/ashmit-singh-gogia/c-hat/internal/routes"
	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	Cfg := config.LoadConfig() // the env vars are initialized inside this method
	config.InitOAuth(Cfg)      // the oauth config is initialized inside this method

	database.ConnectDB(Cfg)
	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	chatRepo := repositories.NewChatRepository(database.DB)
	chatService := services.NewChatService(chatRepo)
	chatHandler := handlers.NewChatHandler(chatService)

	messageRepo := repositories.NewMessageRepository(database.DB)
	messageService := services.NewMessageService(messageRepo)
	messageHandler := handlers.NewMessageHandler(messageService)

	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService, Cfg)

	routes.LoadRoutes(router, userHandler, chatHandler, messageHandler, authHandler)
	fmt.Println("Server running on port", Cfg.PORT)
	if err := router.Run(":" + Cfg.PORT); err != nil {
		panic(err)
	}
}
