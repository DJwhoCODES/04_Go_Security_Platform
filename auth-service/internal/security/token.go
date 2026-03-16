package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GenerateRandomToken(length int) (string, error) {

	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {

	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}
