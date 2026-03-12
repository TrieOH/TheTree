package domain

import "context"

type OAuthProvider interface {
	BuildAuthURL(state, redirectURI string) string
	ExchangeCode(ctx context.Context, code, redirectURI string) (ProviderCredentialData, error)
	Me(ctx context.Context, accessToken string) (int, error)
}
