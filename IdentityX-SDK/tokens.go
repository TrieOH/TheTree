package goauth

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/golang-jwt/jwt/v5"
)

type JWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type TokenService struct {
	client      *Client
	mu          sync.RWMutex
	jwks        *JWKS
	lastUpdated time.Time
}

func (s *TokenService) GetJWKS(ctx context.Context, forceRefresh bool) (*JWKS, error) {
	s.mu.RLock()
	if !forceRefresh && s.jwks != nil {
		defer s.mu.RUnlock()
		return s.jwks, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Re-check after acquiring lock
	if !forceRefresh && s.jwks != nil {
		return s.jwks, nil
	}

	// Cooldown: don't force refresh more than once every 5 minutes
	if forceRefresh && time.Since(s.lastUpdated) < 5*time.Minute {
		if s.jwks != nil {
			return s.jwks, nil
		}
	}

	path := fmt.Sprintf("/projects/%s/.well-known/jwks.json", s.client.projectID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res JWKS
	err = s.client.do(req, &res)
	if err != nil {
		return nil, err
	}

	s.jwks = &res
	s.lastUpdated = time.Now()
	return s.jwks, nil
}

func (s *TokenService) ValidateToken(ctx context.Context, tokenStr string) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fail.New(SDKUnexpectedSigningMethod).WithArgs(token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fail.New(SDKMissingTokenKID)
		}

		// Try with current cache
		jwks, err := s.GetJWKS(ctx, false)
		if err != nil {
			return nil, err
		}

		for _, key := range jwks.Keys {
			if key.Kid == kid {
				return s.decodeKey(key)
			}
		}

		// Kid not found, try refreshing
		jwks, err = s.GetJWKS(ctx, true)
		if err != nil {
			return nil, err
		}

		for _, key := range jwks.Keys {
			if key.Kid == kid {
				return s.decodeKey(key)
			}
		}

		return nil, fail.New(SDKKeyNotInJWKS).WithArgs(kid)
	}

	token, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return nil, fail.New(SDKUnknownErrorID).WithArgs(err.Error())
	}

	return token, nil
}

func (s *TokenService) decodeKey(key JWK) (interface{}, error) {
	if key.Kty != "OKP" || key.Crv != "Ed25519" {
		return nil, fail.New(SDKUnsupportedCurve).WithArgs(key.Kty, key.Crv)
	}

	pubBytes, err := base64.RawURLEncoding.DecodeString(key.X)
	if err != nil {
		return nil, fail.New(SDKKeyDecodeFailed).WithArgs(err.Error())
	}

	if len(pubBytes) != ed25519.PublicKeySize {
		return nil, fail.New(SDKInvalidKeySize).WithArgs(len(pubBytes))
	}

	return ed25519.PublicKey(pubBytes), nil
}
