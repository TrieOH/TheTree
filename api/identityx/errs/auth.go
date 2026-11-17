package errs

import "errors"

// Identity and tokens.
var (
	ErrNoIdentity   = errors.New("identity not found in context")
	ErrInvalidToken = errors.New("invalid token")
)

// API keys.
var (
	ErrInvalidAPIKey = errors.New("invalid api key")
)

// Login, registration and setup.
var (
	ErrTermsNotAccepted = errors.New("you must accept the terms of service to register")
)
