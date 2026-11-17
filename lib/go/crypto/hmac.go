package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HashHMACSHA256 produces a verifiable hash for an interface.
// Uses HMAC-SHA256(value, secret), truncated to 16 hex chars.
// The output is non-guessable without the key and suitable for public /verify/{hash} URLs.
func HashHMACSHA256(value string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(value))
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum)[:16]
}
