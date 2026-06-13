package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(payload []byte, kp *KeyPair) (string, error) {
	sig, err := Sign(kp, payload)
	if err != nil {
		return "", err
	}
	return string(payload) + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

func VerifyToken(tokenStr string, publicKeyPEM string, claims jwt.Claims) (*jwt.Token, error) {
	pubKey, err := parseEd25519Public(publicKeyPEM)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func OpenUnverified(tokenStr string, claims jwt.Claims) (*jwt.Token, error) {
	p := jwt.NewParser()
	token, _, err := p.ParseUnverified(tokenStr, claims)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func parseEd25519Public(pem string) (ed25519.PublicKey, error) {
	pub, err := jwt.ParseEdPublicKeyFromPEM([]byte(pem))
	if err != nil {
		return nil, err
	}
	return pub.(ed25519.PublicKey), nil
}
