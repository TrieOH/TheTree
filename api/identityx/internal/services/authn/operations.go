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
}

func NewOperations(
	actors ports.ActorRepo,
	projects ports.ProjectRepo,
	platformRoles ports.PlatformRolesRepo,
	cryptoKeys ports.CryptoKeysRepo,
	blacklist ports.BlacklistRepo,
	externalIdentities ports.ExternalIdentitiesRepo,
) *Operations {
	return errx.MustProvide(&Operations{
		actors:             actors,
		projects:           projects,
		platformRoles:      platformRoles,
		cryptoKeys:         cryptoKeys,
		blacklist:          blacklist,
		externalIdentities: externalIdentities,
	})
}
