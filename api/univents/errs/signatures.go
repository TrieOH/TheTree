package errs

import "errors"

// Signature requests and their tokens.
var (
	ErrInvalidSignatureToken         = errors.New("invalid or expired token")
	ErrSignatureRequestNotPending    = errors.New("signature request is no longer pending")
	ErrTokenRequestEditionMismatch   = errors.New("token does not match request edition")
	ErrTokenSignatureEditionMismatch = errors.New("token does not match signature edition")
	ErrSignatureRequestHandled       = errors.New("signature request not found or already handled")
	ErrSignatureRequestCompleted     = errors.New("signature request not found or already completed")
)
