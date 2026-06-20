package api_keys

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

const (
	apiKeyAlphabet  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz.-_"
	apiKeySecretLen = 64
	ApiKeyPrefixLen = 12 // chars safe to store/display unhashed
)

type GeneratedAPIKey struct {
	Raw    string // full key, return to caller once, never persisted
	Prefix string // first N chars of the random body, safe to store
	Hash   string // sha256 hex of Raw, stored for verification
}

func GenerateAPIKey(env string) (*GeneratedAPIKey, error) {
	body, err := randomString(apiKeySecretLen)
	if err != nil {
		return nil, fmt.Errorf("generating api key: %w", err)
	}

	raw := fmt.Sprintf("idx_%s_%s", env, body)
	sum := sha256.Sum256([]byte(raw))

	return &GeneratedAPIKey{
		Raw:    raw,
		Prefix: body[:ApiKeyPrefixLen],
		Hash:   hex.EncodeToString(sum[:]),
	}, nil
}

// VerifyAPIKey re-hashes a raw key and checks it against a stored hash.
// Use this on auth, never compare raw strings or decode the hash.
func VerifyAPIKey(raw string, storedHash string) bool {
	sum := sha256.Sum256([]byte(raw))
	computedHash := hex.EncodeToString(sum[:])

	return subtle.ConstantTimeCompare([]byte(computedHash), []byte(storedHash)) == 1
}

func randomString(n int) (string, error) {
	out := make([]byte, n)
	alphabetLen := big.NewInt(int64(len(apiKeyAlphabet)))
	for i := range out {
		idx, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		out[i] = apiKeyAlphabet[idx.Int64()]
	}
	return string(out), nil
}

func StripKeyHeader(rawKey string) (string, error) {
	idx := strings.IndexByte(rawKey, '_')
	if idx == -1 {
		return "", fmt.Errorf("malformed api key")
	}
	idx2 := strings.IndexByte(rawKey[idx+1:], '_')
	if idx2 == -1 {
		return "", fmt.Errorf("malformed api key")
	}
	return rawKey[idx+1+idx2+1:], nil
}
