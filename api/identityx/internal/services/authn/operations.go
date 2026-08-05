package authn

import (
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	actors             ports.ActorRepo
	projects           ports.ProjectRepo
	platformRoles      ports.PlatformRolesRepo
	cryptoKeys         ports.CryptoKeysRepo
	blacklist          ports.BlacklistRepo
	externalIdentities ports.ExternalIdentitiesRepo
	oauthProviders     ports.ProjectOAuthProvidersRepo
	oauthLoginStates   ports.OAuthLoginStatesRepo
	actionTokens       ports.ActionTokenRepo
	emailSender        ports.EmailSender
	hmacSecret         []byte
}

func NewOperations(
	actors ports.ActorRepo,
	projects ports.ProjectRepo,
	platformRoles ports.PlatformRolesRepo,
	cryptoKeys ports.CryptoKeysRepo,
	blacklist ports.BlacklistRepo,
	externalIdentities ports.ExternalIdentitiesRepo,
	oauthProviders ports.ProjectOAuthProvidersRepo,
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
		externalIdentities: externalIdentities,
		oauthProviders:     oauthProviders,
		oauthLoginStates:   oauthLoginStates,
		actionTokens:       actionTokens,
		emailSender:        emailSender,
		hmacSecret:         hmacSecret,
	})
}
