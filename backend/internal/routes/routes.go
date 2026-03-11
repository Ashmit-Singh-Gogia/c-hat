package routes

import (
	"github.com/ashmit-singh-gogia/c-hat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine, handler *handlers.UserHandler) {

	api := router.Group("/api")
	users := api.Group("/users") // creates a user subgroup inside the api group
	users.POST("/", handler.RegisterUser)
}
