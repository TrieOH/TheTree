// Package services aggregates every feature's operations layer. Import this
// package instead of the per-feature subpackages:
//
//	ops := services.NewOperations(r, authzSvc, tokensMgr, hmacSecret, sender)
//	actor, err := ops.Actors.GetByID(ctx, id)
package services

import (
	"time"

	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	"IdentityX/internal/keys"
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
// the api_keys service. actionTokenMgr owns the single-use action-token
// lifecycle: the sender mints through it, authn redeems through it.
// tokensMgr owns the session-token lifecycle; authn crosses it instead of
// touching keys or the blacklist directly. keysMgr owns the Key-lifecycle;
// project creation (projects and organizations) crosses its Ensure seam
// instead of reaching into the crypto-key repo.
func NewOperations(r *repos.Repos, authzSvc *authz.Service, tokensMgr *tokens.Manager, actionTokenMgr *tokens.ActionTokenManager, keysMgr *keys.Manager, hmacSecret string, sender *emails.Sender) *Operations {
	oauthProviders := NewOAuthProviders(
		r.OAuthProviders, r.OAuthProviders, r.Projects, r.ExternalIdentities, r.Actors,
		authzSvc, tokensMgr,
		resty.New().SetTimeout(15*time.Second),
		oauth.Registry,
	)
	return &Operations{
		Actors:         NewActors(r.Actors, r.Projects, authzSvc),
		APIKeys:        NewAPIKeys([]byte(hmacSecret), r.Actors, r.APIKeys, r.Capabilities, r.Projects, authzSvc),
		Authn:          NewAuthn(r.Actors, r.Projects, r.PlatformRoles, tokensMgr, actionTokenMgr, sender),
		Capabilities:   NewCapabilities(r.Actors, r.Capabilities, r.Projects, authzSvc),
		EmailTemplates: NewEmailTemplates(r.EmailTemplates, authzSvc),
		Organizations:  NewOrganizations(r.Projects, r.Actors, r.Organizations, keysMgr, authzSvc),
		OAuthProviders: oauthProviders,
		ProfileSchemas: NewProfileSchemas(r.ProfileSchemas, r.Projects, authzSvc),
		Profiles:       NewProfiles(r.Profiles, r.ProfileSchemas, r.Actors, authzSvc),
		Projects:       NewProjects(r.Projects, r.Actors, keysMgr, authzSvc),
	}
}
