package ports

import (
	"context"
	"payssage/models"
)

type OAuthProvider interface {
	// BuildAuthURL builds the provider consent URL for the OAuth flow. The
	// redirect URI is Payssage's own callback (from config) — the provider
	// owns it, callers never supply it (D7).
	BuildAuthURL(state string) string
	// ExchangeCode exchanges the provider's authorization code for
	// credentials, using Payssage's own configured redirect URI.
	ExchangeCode(ctx context.Context, code string) (models.ProviderCredentialData, error)
}
