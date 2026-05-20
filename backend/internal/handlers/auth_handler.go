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

func (h *AuthHandler) HandleLocalRegister(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := h.authService.RegisterLocalUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
func (h *AuthHandler) HandleLocalLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.authService.LoginLocalUser(req.Email, req.Password)
	if err != nil {
		if err.Error() == "unverified_account" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unverified_account"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	token, err := utils.CreateToken(user.ID, h.cfg.JWT_SECRET)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_token", token, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *AuthHandler) HandleEmailVerification(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}
	err := h.authService.VerifyEmail(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *AuthHandler) HandleResendVerification(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := h.authService.ResendVerificationEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend verification email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Verification email resent successfully"})
}
