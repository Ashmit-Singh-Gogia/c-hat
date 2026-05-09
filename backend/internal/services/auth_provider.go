package services

import "github.com/gin-gonic/gin"

type AuthProvider interface {
	Login(c *gin.Context)
	Logout(c *gin.Context) error
	Callback(c *gin.Context) (*UserDTO, error)
}

type UserDTO struct {
	Provider   string
	ProviderID string
	Email      string
	Username   string
	Avatar     string
}
