package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

// FIXME find a better place for this file that isn't the global utils folder

func GenerateEd25519Keys() (pubKeyPEM, privKeyPEM string, err error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", ErrGeneratingEd25519Key{Cause: err}
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", "", ErrMarshalingPKCS8PrivateKey{Cause: err}
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", "", ErrMarshalingPKIXPublicKey{Cause: err}
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return string(pubPEM), string(privPEM), nil
}

func ParseEd25519PrivateKey(pemStr string) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, ErrInvalidPEMPrivKey{}
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, ErrParsingPKCS8PrivKey{Cause: err}
	}

	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, ErrNotED25519PrivKey{}
	}

	if len(priv) != ed25519.PrivateKeySize {
		return nil, ErrInvalidPrivKeyLength{}
	}

	return priv, nil
}

func ParseEd25519PublicKey(pemStr string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, ErrInvalidPEMPubKey{}
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, ErrParsingPKIXPubKey{Cause: err}
	}

	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, ErrNotED25519PubKey{}
	}

	if len(pub) != ed25519.PublicKeySize {
		return nil, ErrInvalidPubKeyLength{}
	}

	return pub, nil
}

var (
	GoAuthPrivateKey ed25519.PrivateKey
	GoAuthPublicKey  ed25519.PublicKey
)

func LoadEd25519Keys(privatePath, publicPath string) error {
	privBytes, err := os.ReadFile(privatePath)
	if err != nil {
		return err
	}

	pubBytes, err := os.ReadFile(publicPath)
	if err != nil {
		return err
	}

	priv, err := ParseEd25519PrivateKey(string(privBytes))
	if err != nil {
		return err
	}

	pub, err := ParseEd25519PublicKey(string(pubBytes))
	if err != nil {
		return err
	}

	GoAuthPrivateKey = priv
	GoAuthPublicKey = pub

	return nil
}

func PublicKeyToJWK(key ed25519.PublicKey) map[string]any {
	return map[string]any{
		"kty": "OKP",
		"crv": "Ed25519",
		"alg": "EdDSA",
		"use": "sig",
		"kid": "goauth-ed25519",
		"x":   base64.RawURLEncoding.EncodeToString(key),
	}
}

type ErrParseProjectKey struct {
	KeyType string
	Cause   error
}

func (e ErrParseProjectKey) Error() string {
	return "failed to parse project " + e.KeyType + " key"
}

type ErrInvalidPEMPubKey struct{}

func (e ErrInvalidPEMPubKey) Error() string {
	return "invalid PEM public key"
}

type ErrInvalidPEMPrivKey struct{}

func (e ErrInvalidPEMPrivKey) Error() string {
	return "invalid PEM private key"
}

type ErrParsingPKIXPubKey struct {
	Cause error
}

func (e ErrParsingPKIXPubKey) Error() string {
	return "failed to parse PKIX public key"
}

type ErrParsingPKCS8PrivKey struct {
	Cause error
}

func (e ErrParsingPKCS8PrivKey) Error() string {
	return "failed to parse PKCS8 private key"
}

type ErrNotED25519PubKey struct{}

func (e ErrNotED25519PubKey) Error() string {
	return "not an ED25519 public key"
}

type ErrNotED25519PrivKey struct{}

func (e ErrNotED25519PrivKey) Error() string {
	return "not an ED25519 private key"
}

type ErrInvalidPubKeyLength struct{}

func (e ErrInvalidPubKeyLength) Error() string {
	return "invalid public key length"
}

type ErrInvalidPrivKeyLength struct{}

func (e ErrInvalidPrivKeyLength) Error() string {
	return "invalid private key length"
}

type ErrGeneratingEd25519Key struct {
	Cause error
}

func (e ErrGeneratingEd25519Key) Error() string {
	return "failed to generate ed25519 key"
}

type ErrMarshalingPKCS8PrivateKey struct {
	Cause error
}

func (e ErrMarshalingPKCS8PrivateKey) Error() string {
	return "failed to marshal private key"
}

type ErrMarshalingPKIXPublicKey struct {
	Cause error
}

func (e ErrMarshalingPKIXPublicKey) Error() string {
	return "failed to marshal PKIX public key"
}
