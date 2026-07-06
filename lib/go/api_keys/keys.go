package api_keys

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
)

const (
	alphabet      = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	secretLength  = 128
	displayLength = 32
)

type GeneratedAPIKey struct {
	Raw           string
	DisplayPrefix string
	Hash          []byte
}

func GenerateAPIKey(brand string, env string, hmacSecret []byte) (*GeneratedAPIKey, error) {
	secret, err := randomString(secretLength)
	if err != nil {
		return nil, fmt.Errorf("generate secret: %w", err)
	}

	raw := fmt.Sprintf("%s_v1_%s_%s", brand, env, secret)
	hash := hashAPIKey(raw, hmacSecret)

	return &GeneratedAPIKey{
		Raw:           raw,
		DisplayPrefix: fmt.Sprintf("%s_v1_%s_%s", brand, env, secret[:displayLength]),
		Hash:          hash,
	}, nil
}

// VerifyAPIKey re-hashes a raw key and checks it against a stored hash.
// Use this on auth, never compare raw strings or decode the hash.
func VerifyAPIKey(raw string, storedHash []byte, hmacSecret []byte) bool {
	computed := hashAPIKey(raw, hmacSecret)
	return hmac.Equal(computed, storedHash)
}

func hashAPIKey(raw string, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(raw))
	return mac.Sum(nil)
}

func randomString(n int) (string, error) {
	out := make([]byte, n)
	alphabetLen := big.NewInt(int64(len(alphabet)))
	for i := range out {
		idx, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		out[i] = alphabet[idx.Int64()]
	}
	return string(out), nil
}

type APIKey struct {
	Brand         string
	Version       string
	Environment   string
	Secret        string
	DisplayPrefix string
}

func ParseAPIKey(raw string) (*APIKey, error) {
	parts := strings.SplitN(raw, "_", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid api key")
	}

	secret := parts[3]

	if len(secret) < displayLength {
		return nil, fmt.Errorf("invalid api key")
	}

	return &APIKey{
		Brand:         parts[0],
		Version:       parts[1],
		Environment:   parts[2],
		Secret:        secret,
		DisplayPrefix: fmt.Sprintf("%s_%s_%s_%s", parts[0], parts[1], parts[2], secret[:displayLength]),
	}, nil
}
