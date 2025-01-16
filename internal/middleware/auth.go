// internal/middleware/auth.go - api middleware for JWT token authentication
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kajtekajtek/insight-naturae/internal/jwtutils"
)

func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the token from the Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// validate the token
		username, err := jwtutils.ValidateJWT(secret, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// set the username in the context
		c.Set("username", username)
		// continue with the request
		c.Next()
	}
}