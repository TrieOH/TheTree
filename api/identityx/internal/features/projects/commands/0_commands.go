package commands

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Commands struct {
	keys     ports.CryptoKeysRepo
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	authz    *authz.Service
}

func NewCommands(
	keys ports.CryptoKeysRepo,
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	authz *authz.Service,
) *Commands {
	return errx.MustProvide(&Commands{
		keys:     keys,
		projects: projects,
		actors:   actors,
		authz:    authz,
	})
}
