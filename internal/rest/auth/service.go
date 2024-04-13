package auth

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func parseToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the algorithm used for the token
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Load your public key from the environment variable
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("JWT_PUBLIC_KEY")))
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	return token, err
}

func VerifyToken(tokenStr string) error {
	token, err := parseToken(tokenStr)
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
