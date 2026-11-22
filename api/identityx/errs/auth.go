package errs

import "errors"

// Identity and tokens.
var (
	ErrNoIdentity   = errors.New("this operation needs a signed-in caller; send a bearer token or an API key")
	ErrInvalidToken = errors.New("the token couldn't be read or verified; sign in again to get a new one")
	ErrTokenExpired = errors.New("the token has expired; refresh it, or sign in again if the refresh token expired too")
)

// API keys.
var (
	ErrInvalidAPIKey = errors.New("this API key is unknown, revoked or malformed; check the key you sent, or create a new one")
)
