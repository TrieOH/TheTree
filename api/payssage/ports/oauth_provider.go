package ports

import (
	"context"
	"payssage/models"
)

type OAuthProvider interface {
	BuildAuthURL(state, redirectURI string) string
	ExchangeCode(ctx context.Context, code, redirectURI string) (models.ProviderCredentialData, error)
}
