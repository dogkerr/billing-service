package rest

import (
	"dogker/andrenk/billing-service/internal/rest/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
		err := auth.VerifyToken(tokenStr)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
