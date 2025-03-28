package middleware

import (
	"net/http"
	"obs/internal/types"
	"obs/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware protects routes by verifying the JWT from cookies
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow public access to sign-in and sign-up routes
		path := c.Request.URL.Path
		if path == "/signin" || path == "/signup" {
			c.Next()
			return
		}

		// Get the token from the cookie
		token, err := c.Cookie("auth_token")
		if err != nil {
			res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Authentication required"}
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		// Verify the JWT token
		claims, err := utils.VerifyJWT(token)
		if err != nil {
			res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Invalid or expired session"}
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		// Store the user details in the context for later use
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		// Continue to the next handler
		c.Next()
	}
}
