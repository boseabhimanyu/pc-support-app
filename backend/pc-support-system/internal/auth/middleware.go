package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
	c.Abort()
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Cookie("access_token")
		if err != nil {
			unauthorized(c)
			return
		}

		claims, err := ValidateToken(cookie, secret)
		if err != nil {
			unauthorized(c)
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
