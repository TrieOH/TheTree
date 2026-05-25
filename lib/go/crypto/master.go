package crypto

import (
	"encoding/hex"
	"fmt"
	"lib/errx"
	"os"
	"sync"
)

var (
	masterKey     []byte
	masterKeyOnce sync.Once
)

func MasterKey() []byte {
	masterKeyOnce.Do(func() {
		raw := os.Getenv("ENCRYPTION_KEY")
		if raw == "" {
			errx.Exit(nil, "ENCRYPTION_KEY is not set")
		}

		key, err := hex.DecodeString(raw)
		if err != nil {
			errx.Exit(err, "ENCRYPTION_KEY is not valid hex")
		}
		if len(key) != 32 {
			errx.Exit(nil, fmt.Sprintf("ENCRYPTION_KEY must be 32 bytes, got %d", len(key)))
		}

		masterKey = key
	})
	return masterKey
}
