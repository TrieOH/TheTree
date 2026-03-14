package goauth

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessSub struct {
	ID         uuid.UUID        `json:"id"`
	Email      string           `json:"email"`
	ProjectID  *uuid.UUID       `json:"project_id"`
	UserType   string           `json:"user_type"`
	Metadata   *json.RawMessage `json:"metadata"`
	SessionID  uuid.UUID        `json:"session_id"`
	UserAgent  string           `json:"user_agent"`
	UserIP     string           `json:"user_ip"`
	IsVerified bool             `json:"is_verified"`
	FamilyID   uuid.UUID        `json:"family_id"`
	VerifiedAt *time.Time       `json:"verified_at"`
}

type AccessClaims struct {
	Sub AccessSub `json:"sub"`
	jwt.RegisteredClaims
}

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
	cacheTTL    time.Duration
}

// FIXME use something like sqlite to save the token this way if go auth is unavailable momentarily we can still return the token
func (s *TokenService) GetJWKS(ctx context.Context, forceRefresh bool) (*JWKS, error) {
	s.mu.RLock()
	cached := s.jwks
	lastUpdated := s.lastUpdated
	s.mu.RUnlock()

	cacheValid := cached != nil && time.Since(lastUpdated) < s.cacheTTL

	// Cooldown: evita thundering herd em forceRefresh
	inCooldown := time.Since(lastUpdated) < 5*time.Minute

	if cacheValid && !forceRefresh {
		return cached, nil
	}

	if forceRefresh && inCooldown && cached != nil {
		return cached, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Re-check com write lock (outra goroutine pode ter atualizado)
	cacheValid = s.jwks != nil && time.Since(s.lastUpdated) < s.cacheTTL
	inCooldown = time.Since(s.lastUpdated) < 5*time.Minute

	if cacheValid && !forceRefresh {
		return s.jwks, nil
	}
	if forceRefresh && inCooldown && s.jwks != nil {
		return s.jwks, nil
	}

	// Fetch...
	path := fmt.Sprintf("/projects/%s/.well-known/jwks.json", s.client.projectID)
	req, err := s.client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var res JWKS
	if err = s.client.do(req, &res); err != nil {
		if s.jwks != nil {
			return s.jwks, nil // fallback pro cache stale se fetch falhar
		}
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

func (s *TokenService) VerifyAccessToken(ctx context.Context, tokenStr string) (*AccessClaims, error) {
	claims := &AccessClaims{}

	// ----------------------------
	// Parse unverified to get kid
	// ----------------------------
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, _, err := parser.ParseUnverified(tokenStr, claims)
	if err != nil {
		return nil, fail.New(SDKBadRequestID).WithArgs(err.Error())
	}

	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, fail.New(SDKMissingTokenKID)
	}

	// ----------------------------
	// Key resolver (JWKS only)
	// ----------------------------
	keyFunc := func(token *jwt.Token) (interface{}, error) {

		// enforce Ed25519
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fail.New(SDKUnexpectedSigningMethod).
				WithArgs(token.Header["alg"])
		}

		jwks, err := s.GetJWKS(ctx, false)
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

	// ----------------------------
	// Verified parse
	// ----------------------------
	token, err = jwt.ParseWithClaims(tokenStr, claims, keyFunc)
	if err != nil {
		return nil, fail.New(SDKUnauthorizedID).WithArgs(err.Error())
	}

	if !token.Valid {
		return nil, fail.New(SDKUnauthorizedID)
	}

	// ----------------------------
	// Time validation
	// ----------------------------
	now := time.Now()

	if claims.ExpiresAt == nil || now.After(claims.ExpiresAt.Time) {
		return nil, fail.New(SDKUnauthorizedID).WithArgs("token expired")
	}

	if claims.NotBefore != nil && now.Before(claims.NotBefore.Time) {
		return nil, fail.New(SDKUnauthorizedID).
			WithArgs("token not valid yet")
	}

	// ----------------------------
	// Issuer validation
	// ----------------------------
	if claims.Sub.ProjectID != nil {

		expected := "project:" + claims.Sub.ProjectID.String()

		if claims.Issuer != expected {
			return nil, fail.New(SDKUnauthorizedID).
				WithArgs("invalid issuer")
		}

	} else {

		if claims.Issuer != "goauth" {
			return nil, fail.New(SDKUnauthorizedID).
				WithArgs("invalid issuer")
		}
	}

	return claims, nil
}
