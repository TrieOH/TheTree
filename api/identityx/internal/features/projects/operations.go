package projects

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	keys     ports.CryptoKeysRepo
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	authz    *authz.Service
}

func NewOperations(
	keys ports.CryptoKeysRepo,
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		keys:     keys,
		projects: projects,
		actors:   actors,
		authz:    authz,
	})
}
