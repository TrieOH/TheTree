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

func newServer(ops operations, st *setup.State) *server {
	return &server{
		actorsHandlers:         actors.NewHandlers(ops.actors),
		apiKeysHandlers:        api_keys.NewHandlers(ops.apiKeys),
		authnHandlers:          authn.NewHandlers(ops.authn, st),
		capabilitiesHandlers:   capabilities.NewHandlers(ops.capabilities),
		emailTemplatesHandlers: email_templates.NewHandlers(ops.emailTemplates),
		oauthProvidersHandlers: oauth_providers.NewHandlers(ops.oauthProviders),
		organizationsHandlers:  organizations.NewHandlers(ops.organizations),
		profileSchemasHandlers: profile_schemas.NewHandlers(ops.profileSchemas),
		profilesHandlers:       profiles.NewHandlers(ops.profiles),
		projectsHandlers:       projects.NewHandlers(ops.projects),
		tosHandlers:            tos.NewHandlers(ops.tos),
	}
}
