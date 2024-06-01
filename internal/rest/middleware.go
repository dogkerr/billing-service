package rest

import (
	"dogker/andrenk/billing-service/internal/rest/auth"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("missing Authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid Authorization header"))
			return
		}

		tokenString := parts[1]
		publicKey := auth.GetPublicKey()

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if !token.Valid {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}

		// Get the claims
		claims := token.Claims.(*jwt.MapClaims)

		// Extract the userID from the claims
		userID, ok := (*claims)["sub"].(string)
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("missing sub claim"))
			return
		}

		if !c.IsAborted() {
			c.Set("userID", userID)
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
