package idx

import (
	"fmt"
	"time"
)

// --------------------------------------------------------------------------
// Config errors — programmer errors at startup
// --------------------------------------------------------------------------

// ConfigError is returned by NewClient when required configuration is missing
// or invalid.
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("idx: config: %s: %s", e.Field, e.Message)
}

// --------------------------------------------------------------------------
// Token validation errors — the token itself is bad
// --------------------------------------------------------------------------

// InvalidTokenError is returned when a token cannot be parsed or its signature
// is invalid. Cause carries the underlying jwt error.
type InvalidTokenError struct {
	Cause error
}

func (e *InvalidTokenError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("idx: invalid token: %s", e.Cause)
	}
	return "idx: invalid token"
}

func (e *InvalidTokenError) Unwrap() error { return e.Cause }

// TokenExpiredError is returned when the token's exp claim is in the past.
type TokenExpiredError struct {
	ExpiredAt time.Time
}

func (e *TokenExpiredError) Error() string {
	return fmt.Sprintf("idx: token expired at %s", e.ExpiredAt.UTC().Format(time.RFC3339))
}

// TokenNotYetValidError is returned when the token's nbf claim is in the future.
type TokenNotYetValidError struct {
	ValidAt time.Time
}

func (e *TokenNotYetValidError) Error() string {
	return fmt.Sprintf("idx: token not valid until %s", e.ValidAt.UTC().Format(time.RFC3339))
}

// InvalidIssuerError is returned when the token's iss claim does not match
// the expected value for the project.
type InvalidIssuerError struct {
	Got      string
	Expected string
}

func (e *InvalidIssuerError) Error() string {
	return fmt.Sprintf("idx: invalid token issuer: got %q, expected %q", e.Got, e.Expected)
}

// --------------------------------------------------------------------------
// JWKS / key errors — infra or config problems
// --------------------------------------------------------------------------

// KeyNotFoundError is returned when the token's kid header has no matching
// key in the JWKS, even after a forced refresh.
type KeyNotFoundError struct {
	Kid string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("idx: key %q not found in JWKS", e.Kid)
}

// UnsupportedKeyError is returned when a JWKS key has an unexpected type or
// curve. idx only supports OKP/Ed25519.
type UnsupportedKeyError struct {
	Kty string
	Crv string
}

func (e *UnsupportedKeyError) Error() string {
	return fmt.Sprintf("idx: unsupported key type: kty=%s crv=%s (expected OKP/Ed25519)", e.Kty, e.Crv)
}
