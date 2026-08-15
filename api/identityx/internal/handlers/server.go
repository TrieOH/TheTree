// Package handlers implements the generated StrictServerInterface by
// aggregating one subpackage per feature. Each feature package owns the
// methods of its feature; this package only wires them together.
//
// Auth, validation, scope, and error mapping run in the strict middleware
// stack (see internal/app/auth_dispatch.go); the handlers here are pure
// domain logic + fun envelope construction. Platform-vs-project scope is a
// chain concern derived from each operation's x-scope annotation in the
// spec — the Access-check module registers the checker, the resolver
// validates and enforces it, and the parity test pins the surface.
package handlers

import (
	"IdentityX/internal/handlers/actors"
	"IdentityX/internal/handlers/api_keys"
	"IdentityX/internal/handlers/authn"
	"IdentityX/internal/handlers/capabilities"
	"IdentityX/internal/handlers/email_templates"
	"IdentityX/internal/handlers/oauth_providers"
	"IdentityX/internal/handlers/organizations"
	"IdentityX/internal/handlers/profile_schemas"
	"IdentityX/internal/handlers/profiles"
	"IdentityX/internal/handlers/projects"
	"IdentityX/internal/services"
)

// Server implements openapi.StrictServerInterface by embedding every
// feature's Handlers: method promotion satisfies the generated interface
// with no delegation glue. The aliases exist only to give each feature's
// Handlers type a unique embeddable field name — embedding
// *actors.Handlers and *api_keys.Handlers directly would collide on the
// field name "Handlers".
type (
	ActorHandlers         = actors.Handlers
	APIKeyHandlers        = api_keys.Handlers
	AuthnHandlers         = authn.Handlers
	CapabilityHandlers    = capabilities.Handlers
	EmailTemplateHandlers = email_templates.Handlers
	OAuthProviderHandlers = oauth_providers.Handlers
	OrganizationHandlers  = organizations.Handlers
	ProfileSchemaHandlers = profile_schemas.Handlers
	ProfileHandlers       = profiles.Handlers
	ProjectHandlers       = projects.Handlers
)

type Server struct {
	*ActorHandlers
	*APIKeyHandlers
	*AuthnHandlers
	*CapabilityHandlers
	*EmailTemplateHandlers
	*OAuthProviderHandlers
	*OrganizationHandlers
	*ProfileSchemaHandlers
	*ProfileHandlers
	*ProjectHandlers
}

// NewServer wires the per-feature handlers from the services aggregate.
func NewServer(ops *services.Operations) *Server {
	return &Server{
		ActorHandlers:         actors.New(ops.Actors),
		APIKeyHandlers:        api_keys.New(ops.APIKeys),
		AuthnHandlers:         authn.New(ops.Authn),
		CapabilityHandlers:    capabilities.New(ops.Capabilities),
		EmailTemplateHandlers: email_templates.New(ops.EmailTemplates),
		OAuthProviderHandlers: oauth_providers.New(ops.OAuthProviders),
		OrganizationHandlers:  organizations.New(ops.Organizations),
		ProfileSchemaHandlers: profile_schemas.New(ops.ProfileSchemas),
		ProfileHandlers:       profiles.New(ops.Profiles),
		ProjectHandlers:       projects.New(ops.Projects),
	}
}
