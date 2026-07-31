package commands

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Commands struct {
	actors   ports.ActorRepo
	projects ports.ProjectRepo
	authz    *authz.Service
}

func NewCommands(
	actors ports.ActorRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Commands {
	return errx.MustProvide(&Commands{
		actors:   actors,
		projects: projects,
		authz:    authz,
	})
}
