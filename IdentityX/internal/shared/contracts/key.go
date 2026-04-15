package contracts

import (
	"IdentityX/internal/shared/errx"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type Type string
type Usage string
type Status string
type Algorithm string

const (
	TypeGoAuth  Type = "goauth"
	TypeProject Type = "project"
)

const (
	UsageSign   Usage = "sign"
	UsageVerify Usage = "verify"
)

const (
	StatusActive  Status = "active"
	StatusRotated Status = "rotated"
	StatusRevoked Status = "revoked"
)

const (
	AlgEdDSA Algorithm = "EdDSA"
)

type Pair struct {
	ID  uuid.UUID `json:"id"`
	KID string    `json:"kid"`

	// Scope
	KeyType   Type       `json:"key_type"`
	ProjectID *uuid.UUID `json:"project_id"` // nil for goauth

	// Crypto
	Algorithm  Algorithm `json:"algorithm"`
	PublicKey  string    `json:"public_key"`
	PrivateKey []byte    `json:"private_key"` // encrypted

	// Lifecycle
	Usage  Usage  `json:"usage"`
	Status Status `json:"status"`

	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	VerifyExpiresAt time.Time `json:"verify_expires_at"`
}

type PublicKey struct {
	KID       string    `json:"kid"`
	Algorithm Algorithm `json:"algorithm"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func PublicKeyToJWK(k PublicKey) (map[string]any, error) {
	block, _ := pem.Decode([]byte(k.PublicKey))
	if block == nil {
		return nil, fail.New(errx.SYSInvalidPublicKeyPEM).
			WithArgs("failed to decode PEM block")
	}

	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fail.New(errx.SYSPublicKeyParseFailed).
			WithArgs(err.Error())
	}

	pub, ok := pubAny.(ed25519.PublicKey)
	if !ok {
		return nil, fail.New(errx.SYSInvalidPublicKeyType).
			WithArgs("expected Ed25519 public key")
	}

	if len(pub) != ed25519.PublicKeySize {
		return nil, fail.New(errx.SYSInvalidPublicKeyByteSize).
			WithArgs("invalid Ed25519 key length")
	}

	x := base64.RawURLEncoding.EncodeToString(pub)

	return map[string]any{
		"kty": "OKP",
		"crv": "Ed25519",
		"x":   x,
		"alg": "EdDSA",
		"use": "sig",
		"kid": k.KID,
	}, nil
}
