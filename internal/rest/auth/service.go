package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var publicKey *ecdsa.PublicKey

// InitPublicKey initializes the public key used for JWT verification
func InitPublicKey(publicKeyPEM string) error {
	publicKeyBytes := []byte(publicKeyPEM)

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return ErrInvalidPublicKey
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	publicKey = pub.(*ecdsa.PublicKey)
	return nil
}

// GetPublicKey returns the initialized public key
func GetPublicKey() *ecdsa.PublicKey {
	return publicKey
}

// ErrInvalidPublicKey is returned when the provided public key is invalid
var ErrInvalidPublicKey = errors.New("invalid public key")
