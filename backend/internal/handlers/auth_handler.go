package handlers

import (
	"net/http"

	"github.com/ashmit-singh-gogia/c-hat/internal/config"
	"github.com/ashmit-singh-gogia/c-hat/internal/services"
	"github.com/ashmit-singh-gogia/c-hat/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
	authManager *services.AuthManager
	cfg         *config.Config
}

func NewAuthHandler(authService *services.AuthService, authManager *services.AuthManager, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		authManager: authManager,
		cfg:         cfg,
	}
}

func (h *AuthHandler) HandleLogin(c *gin.Context) {
	p := c.Param("provider")
	provider, err := h.authManager.GetProvider(p)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	provider.Login(c)

}

func (h *AuthHandler) HandleCallback(c *gin.Context) {
	p := c.Param("provider")
	provider, err := h.authManager.GetProvider(p)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	userDTO, err := provider.Callback(c)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:5173/login?error=invalid_provider")
		return
	}
	appUser, err := h.authService.ProcessUserLogin(userDTO)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:5173/login?error=user_creation_failed")
		return
	}
	token, err := utils.CreateToken(appUser.ID, h.cfg.JWT_SECRET)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:5173/login?error=token_generation_failed")
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_token", token, 3600*24, "/", "", false, true)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `
    <!DOCTYPE html>
    <html>
        <head><title>Authenticating...</title></head>
        <body>
            <script>
                // Instantly redirect to the frontend via JavaScript
                window.location.href = "http://localhost:5173/";
            </script>
        </body>
    </html>
`)

}

func (h *AuthHandler) HandleLogout(c *gin.Context) {
	p := c.Param("provider")
	provider, err := h.authManager.GetProvider(p)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	provider.Logout(c)
	// Clear the JWT cookie
	c.SetCookie("jwt_token", "", -1, "/", "", false, true)
	// Send user back to frontend
	c.Redirect(http.StatusFound, "http://localhost:5173/login")
}
