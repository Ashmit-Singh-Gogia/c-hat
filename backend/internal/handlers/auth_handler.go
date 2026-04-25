package handlers

import (
	"context"
	"fmt"

	"github.com/ashmit-singh-gogia/c-hat/internal/config"
	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	"github.com/ashmit-singh-gogia/c-hat/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	authService *services.AuthService
	cfg         *config.Config
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "provider", "google")
	c.Request = c.Request.WithContext(ctx)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Println("Goth error:", err)
		c.JSON(500, gin.H{"error": "Failed to authenticate user"})
		return
	}

	appUser, err := h.authService.FindOrCreateUser(user.UserID, user.Email, user.Name, user.AvatarURL)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process user"})
		return
	}

	token, err := utils.CreateToken(appUser.ID, h.cfg.JWT_SECRET)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func NewAuthHandler(authService *services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, cfg: cfg}
}
