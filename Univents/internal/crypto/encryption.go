package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/spf13/viper"
)

var (
	encryptionKey []byte
	initOnce      sync.Once
	errKeyNotSet  = errors.New("ENCRYPTION_KEY is not set or invalid (must be 32 bytes hex)")
)

// InitEncryption loads the encryption key from viper.
// It expects a 64-character hex string (32 bytes).
func InitEncryption() error {
	var err error
	initOnce.Do(func() {
		keyStr := viper.GetString("ENCRYPTION_KEY")
		if keyStr == "" {
			err = errKeyNotSet
			return
		}

		encryptionKey, err = hex.DecodeString(keyStr)
		if err != nil {
			err = fmt.Errorf("failed to decode ENCRYPTION_KEY: %w", err)
			return
		}

		if len(encryptionKey) != 32 {
			err = fmt.Errorf("ENCRYPTION_KEY must be 32 bytes (got %d)", len(encryptionKey))
			return
		}
	})
	return err
}

// Encrypt encrypts data using AES-GCM.
func Encrypt(plaintext []byte) ([]byte, error) {
	if len(encryptionKey) == 0 {
		if err := InitEncryption(); err != nil {
			return nil, err
		}
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts data using AES-GCM.
func Decrypt(ciphertext []byte) ([]byte, error) {
	if len(encryptionKey) == 0 {
		if err := InitEncryption(); err != nil {
			return nil, err
		}
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
