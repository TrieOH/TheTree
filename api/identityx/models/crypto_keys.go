package models

import (
	"encoding/json"
	"lib/crypto"
	"time"

	"github.com/google/uuid"
)

type CryptoKeyType string

const (
	EncryptionCryptoKeyType CryptoKeyType = "encryption"
	SigningCryptoKeyType    CryptoKeyType = "signing"
)

type CryptoKeyStatus string

const (
	CryptoKeyStatusActive   CryptoKeyStatus = "active"
	CryptoKeyStatusRetiring CryptoKeyStatus = "retiring"
	CryptoKeyStatusRetired  CryptoKeyStatus = "retired"
	CryptoKeyStatusRevoked  CryptoKeyStatus = "revoked"
)

type CryptoKey struct {
	ID                  uuid.UUID       `json:"id"`
	ProjectID           *uuid.UUID      `json:"project_id"`
	Type                CryptoKeyType   `json:"type"`
	Status              CryptoKeyStatus `json:"status"`
	PublicKey           string          `json:"public_key"`
	EncryptedPrivateKey string          `json:"-"`
	Algorithm           string          `json:"algorithm"`
	Metadata            json.RawMessage `json:"metadata"`
	Active              bool            `json:"active"`
	CreatedAt           time.Time       `json:"created_at"`
	RotatedAt           *time.Time      `json:"rotated_at"`
	ExpiresAt           *time.Time      `json:"expires_at"`
}

func (c *CryptoKey) ToKeyPair() *crypto.KeyPair {
	return &crypto.KeyPair{
		Public:           c.PublicKey,
		EncryptedPrivate: c.EncryptedPrivateKey,
		Algorithm:        c.Algorithm,
	}
}
