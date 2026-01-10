package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
)

// FIXME: Implement apierr in this file

func GenerateEd25519Keys() (pubKeyPEM, privKeyPEM string, err error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", "", err
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", "", err
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
		return nil, errors.New("invalid PEM private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("not an Ed25519 private key")
	}

	if len(priv) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid private key length")
	}

	return priv, nil
}

func ParseEd25519PublicKey(pemStr string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("invalid PEM public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("not an Ed25519 public key")
	}

	if len(pub) != ed25519.PublicKeySize {
		return nil, errors.New("invalid public key length")
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
