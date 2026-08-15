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
	cryptoKeys    ports.CryptoKeysRepo
	blacklist     ports.BlacklistRepo
	// verifier owns token open/key/verify/blacklist reasoning; refresh and
	// logout cross the same seam as the auth middleware.
	verifier           *tokens.Verifier
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
	cryptoKeys ports.CryptoKeysRepo,
	blacklist ports.BlacklistRepo,
	verifier *tokens.Verifier,
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
		cryptoKeys:         cryptoKeys,
		blacklist:          blacklist,
		verifier:           verifier,
		externalIdentities: externalIdentities,
		oauthProviders:     oauthProviders,
		oauthLoginStates:   oauthLoginStates,
		actionTokens:       actionTokens,
		emailSender:        emailSender,
		hmacSecret:         hmacSecret,
	})
}
