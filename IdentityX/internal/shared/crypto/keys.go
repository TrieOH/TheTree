package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomSecret generates a high-entropy random string of the given byte size,
// encoded using Base64 RawURL.
func GenerateRandomSecret(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashBcryptSecret hashes a secret using bcrypt.
func HashBcryptSecret(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyBcryptSecret verifies a secret against its bcrypt hash.
func VerifyBcryptSecret(hash, secret string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
}
