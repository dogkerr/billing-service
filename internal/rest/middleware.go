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
			fmt.Println("1")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("2")
			c.AbortWithStatus(http.StatusUnauthorized)
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
			fmt.Println("3")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			fmt.Println("4")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Get the claims
		claims := token.Claims.(*jwt.MapClaims)
		fmt.Println("5", claims)

		// Extract the userID from the claims
		userID, ok := (*claims)["sub"].(string)
		if !ok {
			fmt.Println("5b")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !c.IsAborted() {
			fmt.Println("UserID:", userID)
			c.Set("userID", userID)
		}

		c.Next()
	}
}
