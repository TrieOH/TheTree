package domain

import "context"

type OAuthProvider interface {
	BuildAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (ProviderCredentialData, error)
}
