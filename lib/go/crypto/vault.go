package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
)

// Vault seals private key material with the 32-byte master key
// (AES-256-GCM). Build it once at boot from the service config and pass it
// through Deps; nothing in this package reads the environment.
type Vault struct{ key []byte }

// NewVault parses the hex-encoded master key (ENCRYPTION_KEY).
func NewVault(hexKey string) (*Vault, error) {
	if hexKey == "" {
		return nil, errors.New("master key is empty")
	}
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("master key is not valid hex: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("master key must be 32 bytes, got %d", len(key))
	}
	return &Vault{key: key}, nil
}
