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

func (handler *UserHandler) GetMe(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	id, ok := userID.(uint)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid user ID type"})
		return
	}

	user, err := handler.service.GetUserByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}
