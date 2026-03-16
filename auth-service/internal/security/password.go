package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32
)

func generateSalt() ([]byte, error) {

	salt := make([]byte, saltLength)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func HashPassword(password string) (string, error) {

	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s.%s", b64Salt, b64Hash), nil
}

func VerifyPassword(password, encodedHash string) bool {

	var saltB64, hashB64 string

	_, err := fmt.Sscanf(encodedHash, "%s.%s", &saltB64, &hashB64)
	if err != nil {
		return false
	}

	salt, _ := base64.RawStdEncoding.DecodeString(saltB64)
	hash, _ := base64.RawStdEncoding.DecodeString(hashB64)

	comparison := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	return base64.RawStdEncoding.EncodeToString(comparison) ==
		base64.RawStdEncoding.EncodeToString(hash)
}
