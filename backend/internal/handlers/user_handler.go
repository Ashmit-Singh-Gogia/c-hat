package handlers

import (
	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

type CreateUserRequest struct {
	Username string `json:"username"`
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) RegisterUser(c *gin.Context) {
	user := CreateUserRequest{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}
	newUser, err := handler.service.RegisterUser(user.Username)
	if err != nil {
		c.JSON(409, gin.H{
			"Error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"user created": newUser,
		"message":      "success",
	})
}
