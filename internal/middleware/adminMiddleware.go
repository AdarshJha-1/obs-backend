package middleware

import (
	"net/http"
	"obs/internal/types"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the user has the "admin" role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the role from the context
		role, exists := c.Get("role")
		if !exists {
			res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Role not found"}
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		// Check if the role is "admin"
		if role != "admin" {
			res := types.Response{StatusCode: http.StatusForbidden, Success: false, Message: "Access denied, admin role required"}
			c.JSON(http.StatusForbidden, res)
			c.Abort()
			return
		}

		// Continue to the next handler
		c.Next()
	}
}

