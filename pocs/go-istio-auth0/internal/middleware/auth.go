package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(header string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader(header)
		if user == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated identity"})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
