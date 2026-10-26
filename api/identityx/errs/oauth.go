package errs

import "errors"

// The OAuth login flow.
var (
	ErrMissingOAuthCode   = errors.New("missing code")
	ErrMissingOAuthState  = errors.New("missing state")
	ErrOAuthStateNotFound = errors.New("invalid or expired OAuth state; start a new login")
	ErrOAuthStateMismatch = errors.New("invalid OAuth state")
	ErrOAuthStateExpired  = errors.New("OAuth state expired; start a new login")

	ErrOAuthCodeExchange       = errors.New("failed to exchange code")
	ErrOAuthEmptyToken         = errors.New("empty token from provider")
	ErrOAuthUserinfo           = errors.New("userinfo request failed")
	ErrOAuthIncompleteUserinfo = errors.New("incomplete userinfo from provider")
	ErrOAuthGitHubEmail        = errors.New("could not fetch github email")
)

// UnsupportedProviderError is a login provider IdentityX does not know.
type UnsupportedProviderError struct{ Provider string }

func (e *UnsupportedProviderError) Error() string {
	return "unsupported provider: " + e.Provider
}

// ProviderNotConfiguredError is a known provider with no credentials in the
// scope: the platform's env, or the project's provider row.
type ProviderNotConfiguredError struct {
	Provider   string
	ForProject bool
}

func (e *ProviderNotConfiguredError) Error() string {
	if e.ForProject {
		return "provider not configured for this project: " + e.Provider
	}
	return "provider not configured: " + e.Provider
}

// ProviderDisabledError is a sign-up through a provider the project turned
// off; existing identities still log in through it.
type ProviderDisabledError struct{ Provider string }

func (e *ProviderDisabledError) Error() string {
	return "this project has disabled " + e.Provider + " login; contact the project to enable it"
}

// ProviderRemovedError is a login whose provider row was removed while the
// login was in flight (between connect and callback).
type ProviderRemovedError struct{ Provider string }

func (e *ProviderRemovedError) Error() string {
	return "this project has disabled " + e.Provider + " login; contact the project to enable it"
}
