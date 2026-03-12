package domain

import "context"

type OAuthProvider interface {
	BuildAuthURL(state, redirectURI string) string
	ExchangeCode(ctx context.Context, code, redirectURI string) (ProviderCredentialData, error)
	MeID(ctx context.Context, accessToken string) (int, error)
	MeName(ctx context.Context, accessToken string) (string, error)
}
