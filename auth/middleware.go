package auth

import (
	"strings"

	"github.com/OkanUysal/go-response"
	"github.com/gin-gonic/gin"
)

// Middleware validates JWT tokens and sets user info in context
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Check if token is blacklisted (by JTI or family ID)
		if IsTokenBlacklisted(claims.ID) || IsTokenFamilyBlacklisted(claims.FamilyID) {
			response.Unauthorized(c, "Token has been revoked")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetRole retrieves the user role from the context
func GetRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	return role.(string), true
}
