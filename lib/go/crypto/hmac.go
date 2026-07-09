package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

// HashHMACSHA256 produces a verifiable hash for an interface.
// Uses HMAC-SHA256(value, HMAC_SECRET), truncated to 16 hex chars.
// The output is non-guessable without the key and suitable for public /verify/{hash} URLs.
func HashHMACSHA256(value string) string {
	mac := hmac.New(sha256.New, []byte(os.Getenv("HMAC_SECRET")))
	mac.Write([]byte(value))
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum)[:16]
}
