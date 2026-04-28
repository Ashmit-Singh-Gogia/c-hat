package handlers

import (
	"context"
	"net/http"

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
	// FIX: Inject provider context here as well for Chrome
	ctx := context.WithValue(c.Request.Context(), "provider", "google")
	c.Request = c.Request.WithContext(ctx)

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.Redirect(302, "http://localhost:5173/login?error=auth_failed")
		return
	}

	appUser, err := h.authService.FindOrCreateUser(user.UserID, user.Email, user.Name, user.AvatarURL)
	if err != nil {
		c.Redirect(302, "http://localhost:5173/login?error=user_creation_failed")
		return
	}

	token, err := utils.CreateToken(appUser.ID, h.cfg.JWT_SECRET)
	if err != nil {
		c.Redirect(302, "http://localhost:5173/login?error=token_generation_failed")
		return
	}

	// 1. Set SameSite policy to Lax so Safari/Chrome don't block the cookie during redirect
	c.SetSameSite(http.SameSiteLaxMode)

	// 2. Set the HTTP-Only cookie containing the JWT
	// Format: Name, Value, MaxAge, Path, Domain, Secure, HttpOnly
	c.SetCookie("jwt_token", token, 3600*24, "/", "", false, true)

	// 3. Redirect straight back to the React app homepage (no token in the URL!)
	c.Redirect(302, "http://localhost:5173/")
}

func (h *AuthHandler) GoogleLogout(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "provider", "google")
	c.Request = c.Request.WithContext(ctx)
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.SetCookie("jwt_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "http://localhost:5173/login")
}

func NewAuthHandler(authService *services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, cfg: cfg}
}
