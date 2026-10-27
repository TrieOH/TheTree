package errs

import "errors"

// Identity and tokens.
var (
	ErrNoIdentity         = errors.New("identity not found in context")
	ErrPlatformOnly       = errors.New("platform-level authentication required")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrMissingKID         = errors.New("missing kid")
	ErrInvalidKID         = errors.New("invalid kid")
	ErrOutdatedToken      = errors.New("outdated token")
	ErrSigningKeyRevoked  = errors.New("token signing key revoked")
	ErrSigningKeyRetired  = errors.New("token signing key retired")
	ErrTokenRevoked       = errors.New("token has been revoked")
)

// API keys.
var (
	ErrInvalidAPIKey       = errors.New("invalid api key")
	ErrAPIKeyExpired       = errors.New("api key expired")
	ErrInvalidCapabilities = errors.New("invalid capabilities")
)

// Login, registration and setup.
var (
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrPasswordOnlyForHumans = errors.New("only human actors can authenticate via password")
	ErrTermsNotAccepted      = errors.New("you must accept the terms of service to register")
	ErrSetupComplete         = errors.New("setup already complete")
	ErrSetupRequired         = errors.New("please setup IDX first on /auth/setup")
	ErrSetupCheckFailed      = errors.New("could not check setup state")
)

// Action tokens (verify and reset links). Redeem returns exactly one of
// these per failure, so callers branch on the class.
var (
	// ErrActionTokenInvalid covers unparseable, wrongly-purposed, or
	// unknown tokens.
	ErrActionTokenInvalid = errors.New("invalid or expired token")
	// ErrActionTokenExpired covers tokens whose anti-replay record has
	// expired.
	ErrActionTokenExpired = errors.New("token expired")
	// ErrActionTokenUsed covers tokens already consumed, including a
	// concurrent consume.
	ErrActionTokenUsed = errors.New("token already used")
)
