package actors

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	actors   ports.ActorRepo
	projects ports.ProjectRepo
	authz    *authz.Service
}

func NewOperations(
	actors ports.ActorRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		actors:   actors,
		projects: projects,
		authz:    authz,
	})
}
