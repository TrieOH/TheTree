package oauth_providers

import (
	"context"

	"IdentityX/models"
	"lib/crypto"
	"lib/oauth"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// LoginProvider is the resolved credential scope for one OAuth provider:
// the credentials a login flow exchanges codes with, plus whether the
// provider is disabled for new sign-ups.
type LoginProvider struct {
	Provider models.OAuthProvider
	Creds    oauth.Credentials
	Disabled bool
}

// ResolveLoginProvider returns the credentials and enablement for one
// provider in a scope — the platform (projectID nil, env-configured) or a
// project's configured provider row. It is the single place login flows
// learn whether a provider is configured and enabled: connect and callback
// consult this module instead of reaching into the repo, so discovery
// (ListEnabledProviders) and the flow can never drift on enablement.
//
// A missing project row surfaces as NotFound — callers decide how to phrase
// it. Project existence itself is the caller's concern.
func (o *Operations) ResolveLoginProvider(ctx context.Context, provider string, projectID *uuid.UUID) (LoginProvider, error) {
	if projectID == nil {
		creds, ok := oauth.EnvCredentials(provider)
		if !ok {
			return LoginProvider{}, fun.ErrBadRequest("provider not configured: " + provider)
		}
		return LoginProvider{Provider: models.OAuthProvider(provider), Creds: creds}, nil
	}

	row, err := o.providers.GetByProjectAndProvider(ctx, *projectID, models.OAuthProvider(provider))
	if err != nil {
		return LoginProvider{}, err
	}

	secret, err := crypto.DecryptPrivateKey(row.EncryptedClientSecret)
	if err != nil {
		return LoginProvider{}, err
	}

	return LoginProvider{
		Provider: models.OAuthProvider(provider),
		Creds: oauth.Credentials{
			ClientID:     row.ClientID,
			ClientSecret: string(secret),
			RedirectURL:  row.CallbackURL,
		},
		Disabled: !row.Enabled,
	}, nil
}
