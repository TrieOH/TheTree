package app

import (
	"IdentityX/internal/features/actors"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/authn"
	"IdentityX/internal/features/capabilities"
	"IdentityX/internal/features/email_templates"
	"IdentityX/internal/features/oauth_providers"
	"IdentityX/internal/features/organizations"
	"IdentityX/internal/features/profile_schemas"
	"IdentityX/internal/features/profiles"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/tos"
	"IdentityX/internal/openapi"
	"IdentityX/internal/setup"
)

// operations is every feature's operations. Build fills all of it; the
// app tests fill only what their routes reach.
type operations struct {
	actors         *actors.Operations
	apiKeys        *api_keys.Operations
	authn          *authn.Operations
	capabilities   *capabilities.Operations
	emailTemplates *email_templates.Operations
	oauthProviders *oauth_providers.Operations
	organizations  *organizations.Operations
	profileSchemas *profile_schemas.Operations
	profiles       *profiles.Operations
	projects       *projects.Operations
	tos            *tos.Operations
}

// server implements openapi.StrictServerInterface by embedding every
// feature's Handlers: method promotion satisfies the generated interface
// with no delegation glue. The aliases only give each embedded field its
// own name — embedding two *X.Handlers directly would collide on
// "Handlers".
type (
	actorHandlers         = actors.Handlers
	apiKeyHandlers        = api_keys.Handlers
	authnHandlers         = authn.Handlers
	capabilityHandlers    = capabilities.Handlers
	emailTemplateHandlers = email_templates.Handlers
	oauthProviderHandlers = oauth_providers.Handlers
	organizationHandlers  = organizations.Handlers
	profileSchemaHandlers = profile_schemas.Handlers
	profileHandlers       = profiles.Handlers
	projectHandlers       = projects.Handlers
	tosHandlers           = tos.Handlers
)

type server struct {
	*actorHandlers
	*apiKeyHandlers
	*authnHandlers
	*capabilityHandlers
	*emailTemplateHandlers
	*oauthProviderHandlers
	*organizationHandlers
	*profileSchemaHandlers
	*profileHandlers
	*projectHandlers
	*tosHandlers
}

var _ openapi.StrictServerInterface = (*server)(nil)

func newServer(ops operations, st *setup.State) *server {
	return &server{
		actorHandlers:         actors.NewHandlers(ops.actors),
		apiKeyHandlers:        api_keys.NewHandlers(ops.apiKeys),
		authnHandlers:         authn.NewHandlers(ops.authn, st),
		capabilityHandlers:    capabilities.NewHandlers(ops.capabilities),
		emailTemplateHandlers: email_templates.NewHandlers(ops.emailTemplates),
		oauthProviderHandlers: oauth_providers.NewHandlers(ops.oauthProviders),
		organizationHandlers:  organizations.NewHandlers(ops.organizations),
		profileSchemaHandlers: profile_schemas.NewHandlers(ops.profileSchemas),
		profileHandlers:       profiles.NewHandlers(ops.profiles),
		projectHandlers:       projects.NewHandlers(ops.projects),
		tosHandlers:           tos.NewHandlers(ops.tos),
	}
}
