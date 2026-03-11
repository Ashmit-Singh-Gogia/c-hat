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
	repo := repositories.NewUserRepository(database.DB)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)
	routes.LoadRoutes(router, handler)
	fmt.Println("Server running on port", cfg.PORT)
	if err := router.Run(":" + cfg.PORT); err != nil {
		panic(err)
	}
}
