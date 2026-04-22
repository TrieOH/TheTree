package idx

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/TrieOH/sdkkit"
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

func (s *TokenService) GetJWKS(ctx context.Context, forceRefresh bool) (*JWKS, error) {
	s.mu.RLock()
	cached := s.jwks
	lastUpdated := s.lastUpdated
	s.mu.RUnlock()

	cacheValid := cached != nil && time.Since(lastUpdated) < s.cacheTTL
	inCooldown := time.Since(lastUpdated) < 5*time.Minute

	if cacheValid && !forceRefresh {
		return cached, nil
	}
	if forceRefresh && inCooldown && cached != nil {
		return cached, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Re-check with write lock (another goroutine may have updated).
	cacheValid = s.jwks != nil && time.Since(s.lastUpdated) < s.cacheTTL
	inCooldown = time.Since(s.lastUpdated) < 5*time.Minute

	if cacheValid && !forceRefresh {
		return s.jwks, nil
	}
	if forceRefresh && inCooldown && s.jwks != nil {
		return s.jwks, nil
	}

	var res JWKS
	path := fmt.Sprintf("/.well-known/jwks.json?project_id=%s", s.client.projectID)
	if err := s.client.DoRequestRaw(ctx, "GET", path, nil, &res); err != nil {
		if s.jwks != nil {
			return s.jwks, nil // stale cache fallback on network error
		}
		return nil, err // *sdkkit.SDKError or *sdkkit.APIError — both appropriate here
	}

	s.jwks = &res
	s.lastUpdated = time.Now()
	return s.jwks, nil
}

// ValidateToken parses and verifies the signature of a raw JWT string,
// returning the raw *jwt.Token. Prefer VerifyAccessToken for full claim
// validation.
func (s *TokenService) ValidateToken(ctx context.Context, tokenStr string) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, &InvalidTokenError{Cause: fmt.Errorf("unexpected signing method: %T", token.Method)}
		}

		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, &InvalidTokenError{Cause: fmt.Errorf("missing kid header")}
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

		// Not in cache — try a forced refresh once.
		jwks, err = s.GetJWKS(ctx, true)
		if err != nil {
			return nil, err
		}
		for _, key := range jwks.Keys {
			if key.Kid == kid {
				return s.decodeKey(key)
			}
		}

		return nil, &KeyNotFoundError{Kid: kid}
	}

	token, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return nil, &InvalidTokenError{Cause: err}
	}
	return token, nil
}

func (s *TokenService) decodeKey(key JWK) (interface{}, error) {
	if key.Kty != "OKP" || key.Crv != "Ed25519" {
		return nil, &UnsupportedKeyError{Kty: key.Kty, Crv: key.Crv}
	}

	pubBytes, err := base64.RawURLEncoding.DecodeString(key.X)
	if err != nil {
		return nil, &sdkkit.SDKError{Op: "decode jwks key", Cause: err}
	}

	if len(pubBytes) != ed25519.PublicKeySize {
		return nil, &sdkkit.SDKError{
			Op:    "decode jwks key",
			Cause: fmt.Errorf("invalid key size: %s bytes", strconv.Itoa(len(pubBytes))),
		}
	}

	return ed25519.PublicKey(pubBytes), nil
}

// VerifyAccessToken fully validates a raw access token string: signature,
// expiry, nbf, and issuer. Returns the parsed AccessClaims on success.
func (s *TokenService) VerifyAccessToken(ctx context.Context, tokenStr string) (*AccessClaims, error) {
	claims := &AccessClaims{}

	// Parse unverified first to extract kid so we can fetch the right key
	// without validating claims we haven't checked yet.
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(tokenStr, claims)
	if err != nil {
		return nil, &InvalidTokenError{Cause: err}
	}

	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, &InvalidTokenError{Cause: fmt.Errorf("missing kid header")}
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, &InvalidTokenError{Cause: fmt.Errorf("unexpected signing method: %T", token.Method)}
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
		return nil, &KeyNotFoundError{Kid: kid}
	}

	token, err = jwt.ParseWithClaims(tokenStr, claims, keyFunc)
	if err != nil {
		return nil, &InvalidTokenError{Cause: err}
	}
	if !token.Valid {
		return nil, &InvalidTokenError{}
	}

	// Time validation.
	now := time.Now()
	if claims.ExpiresAt == nil || now.After(claims.ExpiresAt.Time) {
		expAt := time.Time{}
		if claims.ExpiresAt != nil {
			expAt = claims.ExpiresAt.Time
		}
		return nil, &TokenExpiredError{ExpiredAt: expAt}
	}
	if claims.NotBefore != nil && now.Before(claims.NotBefore.Time) {
		return nil, &TokenNotYetValidError{ValidAt: claims.NotBefore.Time}
	}

	// Issuer validation.
	expectedIssuer := "IdentityX"
	if claims.Sub.ProjectID != nil {
		expectedIssuer = claims.Sub.ProjectID.String()
	}
	if claims.Issuer != expectedIssuer {
		return nil, &InvalidIssuerError{Got: claims.Issuer, Expected: expectedIssuer}
	}

	return claims, nil
}
