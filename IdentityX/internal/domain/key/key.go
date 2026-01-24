package key

import (
	"encoding/base64"
	"strings"
	"time"

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
	AlgEd25519 Algorithm = "Ed25519"
)

type Pair struct {
	ID  uuid.UUID
	KID string

	// Scope
	KeyType   Type
	ProjectID *uuid.UUID // nil for goauth

	// Crypto
	Algorithm  Algorithm
	PublicKey  string
	PrivateKey []byte // encrypted

	// Lifecycle
	Usage  Usage
	Status Status

	CreatedAt time.Time
	ExpiresAt time.Time
}

type PublicKey struct {
	KID       string
	Algorithm Algorithm
	PublicKey string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// PublicKeyToJWK converts a PublicKey into JWKS-compatible map
func PublicKeyToJWK(k PublicKey) map[string]any {
	// Strip PEM headers if needed
	pubBytes := decodePEM(k.PublicKey)

	// Base64URL encode
	x := base64.RawURLEncoding.EncodeToString(pubBytes)

	return map[string]any{
		"kty": "OKP",       // Key Type
		"crv": "Ed25519",   // Curve
		"x":   x,           // Base64URL public key
		"alg": k.Algorithm, // Algorithm (EdDSA)
		"use": "sig",       // Usage: signature
		"kid": k.KID,       // Key ID
	}
}

// decodePEM strips PEM headers and returns raw key bytes
func decodePEM(pem string) []byte {
	// Remove header/footer if present
	pem = strings.ReplaceAll(pem, "-----BEGIN PUBLIC KEY-----", "")
	pem = strings.ReplaceAll(pem, "-----END PUBLIC KEY-----", "")
	pem = strings.ReplaceAll(pem, "\n", "")

	// Decode base64
	raw, _ := base64.StdEncoding.DecodeString(pem)
	return raw
}
