// Package repos aggregates every feature's repository layer. Import this
// package instead of the per-feature subpackages:
//
//	r := repos.New(q)
//	actor, err := r.Actors.GetByID(ctx, id)
package repos

import (
	"IdentityX/internal/sqlc"

	"IdentityX/internal/repos/action_tokens"
	"IdentityX/internal/repos/actors"
	"IdentityX/internal/repos/api_keys"
	"IdentityX/internal/repos/authn"
	"IdentityX/internal/repos/blacklist"
	"IdentityX/internal/repos/capabilities"
	"IdentityX/internal/repos/crypto_keys"
	"IdentityX/internal/repos/email_templates"
	"IdentityX/internal/repos/oauth_providers"
	"IdentityX/internal/repos/organizations"
	"IdentityX/internal/repos/platform_roles"
	"IdentityX/internal/repos/profile_schemas"
	"IdentityX/internal/repos/profiles"
	"IdentityX/internal/repos/projects"
)

// Type and constructor aliases for each feature's repo package.
type (
	Actors         = actors.Repo
	APIKeys        = api_keys.Repo
	Authn          = authn.Repo
	ActionTokens   = action_tokens.Repo
	Blacklist      = blacklist.Repo
	Capabilities   = capabilities.Repo
	CryptoKeys     = crypto_keys.Repo
	EmailTemplates = email_templates.Repo
	Organizations  = organizations.Repo
	OAuthProviders = oauth_providers.Repo
	PlatformRoles  = platform_roles.Repo
	ProfileSchemas = profile_schemas.Repo
	Profiles       = profiles.Repo
	Projects       = projects.Repo
)

var (
	NewActors         = actors.NewRepo
	NewAPIKeys        = api_keys.NewRepo
	NewAuthn          = authn.NewRepo
	NewActionTokens   = action_tokens.NewRepo
	NewBlacklist      = blacklist.NewRepo
	NewCapabilities   = capabilities.NewRepo
	NewCryptoKeys     = crypto_keys.NewRepo
	NewEmailTemplates = email_templates.NewRepo
	NewOrganizations  = organizations.NewRepo
	NewOAuthProviders = oauth_providers.NewRepo
	NewPlatformRoles  = platform_roles.NewRepo
	NewProfileSchemas = profile_schemas.NewSchemaRepo
	NewProfiles       = profiles.NewProfileRepo
	NewProjects       = projects.NewRepo
)

// Repos is the aggregate of every feature repo, constructed once at startup.
type Repos struct {
	Actors         *Actors
	APIKeys        *APIKeys
	Authn          *Authn
	ActionTokens   *ActionTokens
	Blacklist      *Blacklist
	Capabilities   *Capabilities
	CryptoKeys     *CryptoKeys
	EmailTemplates *EmailTemplates
	Organizations  *Organizations
	OAuthProviders *OAuthProviders
	PlatformRoles  *PlatformRoles
	ProfileSchemas *ProfileSchemas
	Profiles       *Profiles
	Projects       *Projects
	// ExternalIdentities is provided by the authn repo.
	ExternalIdentities *Authn
}

// New constructs every feature repo from the shared query handle.
func New(q *sqlc.Queries) *Repos {
	return &Repos{
		Actors:             NewActors(q),
		APIKeys:            NewAPIKeys(q),
		Authn:              NewAuthn(q),
		ActionTokens:       NewActionTokens(q),
		Blacklist:          NewBlacklist(q),
		Capabilities:       NewCapabilities(q),
		CryptoKeys:         NewCryptoKeys(q),
		EmailTemplates:     NewEmailTemplates(q),
		Organizations:      NewOrganizations(q),
		OAuthProviders:     NewOAuthProviders(q),
		PlatformRoles:      NewPlatformRoles(q),
		ProfileSchemas:     NewProfileSchemas(q),
		Profiles:           NewProfiles(q),
		Projects:           NewProjects(q),
		ExternalIdentities: NewAuthn(q),
	}
}
