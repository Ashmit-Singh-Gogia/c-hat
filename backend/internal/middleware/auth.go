package middleware

import (
	"net/http"
	"strings"

	"github.com/ashmit-singh-gogia/c-hat/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read JWT from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			return
		}

		// Expect "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		correctedClaims := claims.(*utils.Claims)
		// Inject userID into context for downstream handlers
		c.Set("userID", correctedClaims.UserID)
		c.Next()
	}
}
