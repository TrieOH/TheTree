package idx

import (
	"errors"
	"fmt"
	"net/http"
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
	return "idx: token expired at " + e.ExpiredAt.UTC().Format(time.RFC3339)
}

// TokenNotYetValidError is returned when the token's nbf claim is in the future.
type TokenNotYetValidError struct {
	ValidAt time.Time
}

func (e *TokenNotYetValidError) Error() string {
	return "idx: token not valid until " + e.ValidAt.UTC().Format(time.RFC3339)
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

// --------------------------------------------------------------------------
// API errors — IdentityX answered with a 4xx/5xx
// --------------------------------------------------------------------------

// APIError is returned when IdentityX responds with a failure status.
// StatusCode is the HTTP status; the rest is the decoded error body, which
// may be empty when the body was not the usual JSON error shape.
type APIError struct {
	StatusCode int      `json:"-"`
	Module     string   `json:"module"`
	Message    string   `json:"message"`
	Code       int      `json:"code"`
	Trace      []string `json:"trace"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("idx: api error %d: %s", e.StatusCode, e.Message)
}

// DecodeKeyError is returned when a JWKS key's public key material is not a
// valid base64url Ed25519 key.
type DecodeKeyError struct {
	Kid   string
	Cause error
}

func (e *DecodeKeyError) Error() string {
	return fmt.Sprintf("idx: decode jwks key %q: %s", e.Kid, e.Cause)
}

func (e *DecodeKeyError) Unwrap() error { return e.Cause }

func IsNotFound(err error) bool      { return hasStatus(err, http.StatusNotFound) }
func IsUnauthorized(err error) bool  { return hasStatus(err, http.StatusUnauthorized) }
func IsForbidden(err error) bool     { return hasStatus(err, http.StatusForbidden) }
func IsConflict(err error) bool      { return hasStatus(err, http.StatusConflict) }
func IsUnprocessable(err error) bool { return hasStatus(err, http.StatusUnprocessableEntity) }

// StatusCode returns the HTTP status behind an *APIError, or 0 for any other
// error.
func StatusCode(err error) int {
	if e, ok := errors.AsType[*APIError](err); ok {
		return e.StatusCode
	}
	return 0
}

func hasStatus(err error, code int) bool {
	return StatusCode(err) == code
}
