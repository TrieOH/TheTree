// Package oauth_providers owns everything about OAuth providers: the
// configured/enabled policy (CRUD, discovery), provider credential
// resolution, and the login flow (connect, callback, one-time login state,
// userinfo fetch, identity linking) at one interface. authn crosses it for
// the flow; handlers cross it for policy. Provider metadata and the HTTP
// client are injected at construction, so tests point both at a local
// server and never mutate a global registry.
package oauth_providers

import (
	"IdentityX/internal/authz"
	"IdentityX/internal/tokens"
	"IdentityX/ports"
	"lib/errx"
	"lib/oauth"

	"resty.dev/v3"
)

type Operations struct {
	providers   ports.ProjectOAuthProvidersRepo
	loginStates ports.OAuthLoginStatesRepo
	projects    ports.ProjectRepo
	external    ports.ExternalIdentitiesRepo
	actors      ports.ActorRepo
	authz       *authz.Service
	// tokens owns the JWT lifecycle; the callback crosses it to mint the
	// session pair for a signed-in actor.
	tokens *tokens.Manager
	// client and meta are the HTTP seam: the token exchange and userinfo
	// calls go through the injected resty client, and provider metadata
	// (endpoints, userinfo URL) comes from the injected registry.
	client *resty.Client
	meta   map[string]oauth.Provider
}

func NewOperations(
	providers ports.ProjectOAuthProvidersRepo,
	loginStates ports.OAuthLoginStatesRepo,
	projects ports.ProjectRepo,
	external ports.ExternalIdentitiesRepo,
	actors ports.ActorRepo,
	authz *authz.Service,
	tokensMgr *tokens.Manager,
	client *resty.Client,
	meta map[string]oauth.Provider,
) *Operations {
	return errx.MustProvide(&Operations{
		providers:   providers,
		loginStates: loginStates,
		projects:    projects,
		external:    external,
		actors:      actors,
		authz:       authz,
		tokens:      tokensMgr,
		client:      client,
		meta:        meta,
	})
}
