package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"errors"
)

func GenerateEd25519Keys() (pubKeyPEM, privKeyPEM string, err error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: priv,
	})
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pub,
	})

	return string(pubPEM), string(privPEM), nil
}

func ParseEd25519PrivateKey(pemStr string) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("invalid private key")
	}
	return ed25519.PrivateKey(block.Bytes), nil
}

func ParseEd25519PublicKey(pemStr string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("invalid public key")
	}
	return ed25519.PublicKey(block.Bytes), nil
}
