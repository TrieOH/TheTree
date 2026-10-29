// Package testrepos is for tests only: every feature's repo over one query
// handle, for the integration tests that seed and assert through several
// features. Build wires repos one by one (ADR-0015); nothing outside tests
// imports this.
package testrepos

import (
	"IdentityX/internal/features/action_tokens"
	"IdentityX/internal/features/actors"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/authn"
	"IdentityX/internal/features/blacklist"
	"IdentityX/internal/features/capabilities"
	"IdentityX/internal/features/crypto_keys"
	"IdentityX/internal/features/email_templates"
	"IdentityX/internal/features/oauth_providers"
	"IdentityX/internal/features/organizations"
	"IdentityX/internal/features/platform_roles"
	"IdentityX/internal/features/profile_schemas"
	"IdentityX/internal/features/profiles"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/tos"
	"IdentityX/internal/sqlc"
)

type Repos struct {
	Actors             *actors.Repo
	APIKeys            *api_keys.Repo
	Authn              *authn.Repo
	ActionTokens       *action_tokens.Repo
	Blacklist          *blacklist.Repo
	Capabilities       *capabilities.Repo
	CryptoKeys         *crypto_keys.Repo
	EmailTemplates     *email_templates.Repo
	Organizations      *organizations.Repo
	OAuthProviders     *oauth_providers.Repo
	PlatformRoles      *platform_roles.Repo
	ProfileSchemas     *profile_schemas.Repo
	Profiles           *profiles.Repo
	Projects           *projects.Repo
	Tos                *tos.Repo
	ExternalIdentities *authn.Repo
}

func New(q *sqlc.Queries) *Repos {
	authnRepo := authn.NewRepo(q)
	return &Repos{
		Actors:             actors.NewRepo(q),
		APIKeys:            api_keys.NewRepo(q),
		Authn:              authnRepo,
		ActionTokens:       action_tokens.NewRepo(q),
		Blacklist:          blacklist.NewRepo(q),
		Capabilities:       capabilities.NewRepo(q),
		CryptoKeys:         crypto_keys.NewRepo(q),
		EmailTemplates:     email_templates.NewRepo(q),
		Organizations:      organizations.NewRepo(q),
		OAuthProviders:     oauth_providers.NewRepo(q),
		PlatformRoles:      platform_roles.NewRepo(q),
		ProfileSchemas:     profile_schemas.NewRepo(q),
		Profiles:           profiles.NewRepo(q),
		Projects:           projects.NewRepo(q),
		Tos:                tos.NewRepo(q),
		ExternalIdentities: authnRepo,
	}
}
