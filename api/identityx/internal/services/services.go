// Package services aggregates every feature's operations layer. Import this
// package instead of the per-feature subpackages:
//
//	ops := services.NewOperations(r, authzSvc, tokensMgr, hmacSecret, sender)
//	actor, err := ops.Actors.GetByID(ctx, id)
package services

import (
	"time"

	"IdentityX/internal/repos"
	"IdentityX/internal/services/actors"
	"IdentityX/internal/services/api_keys"
	"IdentityX/internal/services/authn"
	"IdentityX/internal/services/capabilities"
	"IdentityX/internal/services/email_templates"
	"IdentityX/internal/services/oauth_providers"
	"IdentityX/internal/services/organizations"
	"IdentityX/internal/services/profile_schemas"
	"IdentityX/internal/services/profiles"
	"IdentityX/internal/services/projects"
	"IdentityX/internal/tokens"

	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	"lib/oauth"

	"resty.dev/v3"
)

// Type and constructor aliases for each feature's operations package.
type (
	Actors         = actors.Operations
	APIKeys        = api_keys.Operations
	Authn          = authn.Operations
	Capabilities   = capabilities.Operations
	EmailTemplates = email_templates.Operations
	OAuthProviders = oauth_providers.Operations
	Organizations  = organizations.Operations
	ProfileSchemas = profile_schemas.Operations
	Profiles       = profiles.Operations
	Projects       = projects.Operations
)

var (
	NewActors         = actors.NewOperations
	NewAPIKeys        = api_keys.NewOperations
	NewAuthn          = authn.NewOperations
	NewCapabilities   = capabilities.NewOperations
	NewEmailTemplates = email_templates.NewOperations
	NewOAuthProviders = oauth_providers.NewOperations
	NewOrganizations  = organizations.NewOperations
	NewProfileSchemas = profile_schemas.NewOperations
	NewProfiles       = profiles.NewOperations
	NewProjects       = projects.NewOperations
)

// Operations is the aggregate of every feature's operations, constructed
// once at startup and consumed by the HTTP handlers.
type Operations struct {
	Actors         *Actors
	APIKeys        *APIKeys
	Authn          *Authn
	Capabilities   *Capabilities
	EmailTemplates *EmailTemplates
	OAuthProviders *OAuthProviders
	Organizations  *Organizations
	ProfileSchemas *ProfileSchemas
	Profiles       *Profiles
	Projects       *Projects
}

// NewOperations wires every feature's operations from the shared repos.
// hmacSecret is the API-key signing secret (app config), passed through to
// the api_keys and authn services. sender dispatches the async verify/reset
// emails minted by authn. tokensMgr owns the token lifecycle; authn crosses
// it instead of touching keys or the blacklist directly.
func NewOperations(r *repos.Repos, authzSvc *authz.Service, tokensMgr *tokens.Manager, hmacSecret string, sender *emails.Sender) *Operations {
	oauthProviders := NewOAuthProviders(
		r.OAuthProviders, r.OAuthProviders, r.Projects, r.ExternalIdentities, r.Actors,
		authzSvc, tokensMgr,
		resty.New().SetTimeout(15*time.Second),
		oauth.Registry,
	)
	return &Operations{
		Actors:         NewActors(r.Actors, r.Projects, authzSvc),
		APIKeys:        NewAPIKeys([]byte(hmacSecret), r.Actors, r.APIKeys, r.Capabilities, r.Projects, authzSvc),
		Authn:          NewAuthn(r.Actors, r.Projects, r.PlatformRoles, tokensMgr, r.ActionTokens, sender, []byte(hmacSecret)),
		Capabilities:   NewCapabilities(r.Actors, r.Capabilities, r.Projects, authzSvc),
		EmailTemplates: NewEmailTemplates(r.EmailTemplates, authzSvc),
		Organizations:  NewOrganizations(r.Projects, r.Actors, r.Organizations, authzSvc),
		OAuthProviders: oauthProviders,
		ProfileSchemas: NewProfileSchemas(r.ProfileSchemas, r.Projects, authzSvc),
		Profiles:       NewProfiles(r.Profiles, r.ProfileSchemas, r.Actors, authzSvc),
		Projects:       NewProjects(r.CryptoKeys, r.Projects, r.Actors, authzSvc),
	}
}
