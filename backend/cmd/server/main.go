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

	cfg := config.LoadConfig()
	database.ConnectDB(cfg)
	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	chatRepo := repositories.NewChatRepository(database.DB)
	chatService := services.NewChatService(chatRepo)
	chatHandler := handlers.NewChatHandler(chatService)

	routes.LoadRoutes(router, userHandler, chatHandler)
	fmt.Println("Server running on port", cfg.PORT)
	if err := router.Run(":" + cfg.PORT); err != nil {
		panic(err)
	}
}
