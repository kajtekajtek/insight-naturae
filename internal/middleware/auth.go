// internal/middleware/auth.go - api middleware for JWT token authentication
package middleware

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kajtekajtek/insight-naturae/internal/jwtutils"
)

func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			c.Abort()
			return
		}

		// extract the token from the header
		var token string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization Header"})
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
		c.Request.Header.Set("Username", username)
		// continue with the request
		c.Next()
	}
}