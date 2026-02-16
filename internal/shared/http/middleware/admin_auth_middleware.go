package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminAuthMiddleware struct {
	apiKey string
}

func NewAdminAuthMiddleware(apiKey string) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{apiKey: apiKey}
}

func (m *AdminAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Admin-Key")
		if key == "" || subtle.ConstantTimeCompare([]byte(key), []byte(m.apiKey)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}
