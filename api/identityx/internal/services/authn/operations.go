package authn

import (
	"IdentityX/internal/services/oauth_providers"
	"IdentityX/internal/tokens"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	actors        ports.ActorRepo
	projects      ports.ProjectRepo
	platformRoles ports.PlatformRolesRepo
	// tokens owns the token lifecycle (verify/mint/rotate/revoke); login,
	// refresh, logout, and the OAuth callback cross it instead of touching
	// keys, blacklist, or token claims directly.
	tokens             *tokens.Manager
	externalIdentities ports.ExternalIdentitiesRepo
	// oauthProviders owns provider policy (configured/enabled); the OAuth
	// flow consults it instead of the repo.
	oauthProviders   *oauth_providers.Operations
	oauthLoginStates ports.OAuthLoginStatesRepo
	actionTokens     ports.ActionTokenRepo
	emailSender      ports.EmailSender
	hmacSecret       []byte
}

func NewOperations(
	actors ports.ActorRepo,
	projects ports.ProjectRepo,
	platformRoles ports.PlatformRolesRepo,
	tokensMgr *tokens.Manager,
	externalIdentities ports.ExternalIdentitiesRepo,
	oauthProviders *oauth_providers.Operations,
	oauthLoginStates ports.OAuthLoginStatesRepo,
	actionTokens ports.ActionTokenRepo,
	emailSender ports.EmailSender,
	hmacSecret []byte,
) *Operations {
	return errx.MustProvide(&Operations{
		actors:             actors,
		projects:           projects,
		platformRoles:      platformRoles,
		tokens:             tokensMgr,
		externalIdentities: externalIdentities,
		oauthProviders:     oauthProviders,
		oauthLoginStates:   oauthLoginStates,
		actionTokens:       actionTokens,
		emailSender:        emailSender,
		hmacSecret:         hmacSecret,
	})
}
