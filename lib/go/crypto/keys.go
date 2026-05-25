package crypto

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/ed25519"
)

type KeyPair struct {
	Public           string
	EncryptedPrivate string
	Algorithm        string
}

func GenerateKeyPair(keyType string) (*KeyPair, error) {
	switch keyType {
	case "signing":
		return generateEd25519()
	case "encryption":
		return generateRSA()
	default:
		return nil, fmt.Errorf("unknown key type: %s", keyType)
	}
}

func generateEd25519() (*KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	encryptedPriv, err := EncryptPrivateKey(privBytes)
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		Public:           string(pubPEM),
		EncryptedPrivate: encryptedPriv,
		Algorithm:        "Ed25519",
	}, nil
}

func generateRSA() (*KeyPair, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}

	encryptedPriv, err := EncryptPrivateKey(privBytes)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Public:           string(pubPEM),
		EncryptedPrivate: encryptedPriv,
		Algorithm:        "RSA-4096",
	}, nil
}

// EncryptPrivateKey encrypts raw key bytes using AES-256-GCM with the master key from env.
// Output format: hex(nonce) + ":" + hex(ciphertext)
func EncryptPrivateKey(privBytes []byte) (string, error) {
	block, err := aes.NewCipher(MasterKey())
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, privBytes, nil)

	return hex.EncodeToString(nonce) + ":" + hex.EncodeToString(ciphertext), nil
}

func DecryptPrivateKey(encrypted string) ([]byte, error) {
	parts := strings.SplitN(encrypted, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid encrypted key format")
	}
	nonce, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	ciphertext, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(MasterKey())
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func Sign(kp *KeyPair, payload []byte) ([]byte, error) {
	privBytes, err := DecryptPrivateKey(kp.EncryptedPrivate)
	if err != nil {
		return nil, err
	}
	priv, err := x509.ParsePKCS8PrivateKey(privBytes)
	if err != nil {
		return nil, err
	}
	switch key := priv.(type) {
	case ed25519.PrivateKey:
		return ed25519.Sign(key, payload), nil
	case *rsa.PrivateKey:
		hash := crypto.SHA256
		h := hash.New()
		h.Write(payload)
		return rsa.SignPKCS1v15(rand.Reader, key, hash, h.Sum(nil))
	default:
		return nil, fmt.Errorf("unsupported key type: %T", priv)
	}
}
