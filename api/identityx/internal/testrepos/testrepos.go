// Package testrepos is for tests only: every feature's repo over one query
// handle, for the integration tests that seed and assert through several
// features. Build wires repos one by one (ADR-0015); nothing outside tests
// imports this.
package testrepos

import (
	action_tokens "IdentityX/internal/features/action_tokens/repo"
	actors "IdentityX/internal/features/actors/repo"
	api_keys "IdentityX/internal/features/api_keys/repo"
	authn "IdentityX/internal/features/authn/repo"
	blacklist "IdentityX/internal/features/blacklist/repo"
	capabilities "IdentityX/internal/features/capabilities/repo"
	crypto_keys "IdentityX/internal/features/crypto_keys/repo"
	email_templates "IdentityX/internal/features/email_templates/repo"
	oauth_providers "IdentityX/internal/features/oauth_providers/repo"
	organizations "IdentityX/internal/features/organizations/repo"
	platform_roles "IdentityX/internal/features/platform_roles/repo"
	profile_schemas "IdentityX/internal/features/profile_schemas/repo"
	profiles "IdentityX/internal/features/profiles/repo"
	projects "IdentityX/internal/features/projects/repo"
	tos "IdentityX/internal/features/tos/repo"
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
	authnRepo := authn.New(q)
	return &Repos{
		Actors:             actors.New(q),
		APIKeys:            api_keys.New(q),
		Authn:              authnRepo,
		ActionTokens:       action_tokens.New(q),
		Blacklist:          blacklist.New(q),
		Capabilities:       capabilities.New(q),
		CryptoKeys:         crypto_keys.New(q),
		EmailTemplates:     email_templates.New(q),
		Organizations:      organizations.New(q),
		OAuthProviders:     oauth_providers.New(q),
		PlatformRoles:      platform_roles.New(q),
		ProfileSchemas:     profile_schemas.New(q),
		Profiles:           profiles.New(q),
		Projects:           projects.New(q),
		Tos:                tos.New(q),
		ExternalIdentities: authnRepo,
	}
}
