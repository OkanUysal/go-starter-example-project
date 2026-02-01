package auth

import (
	"github.com/OkanUysal/go-response"
	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := GetRole(c)
		if !exists {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		if role != "ADMIN" {
			response.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
