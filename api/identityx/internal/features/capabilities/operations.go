package capabilities

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	actors       ports.ActorRepo
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
	authz        *authz.Service
}

func NewOperations(
	actors ports.ActorRepo,
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		actors:       actors,
		capabilities: capabilities,
		projects:     projects,
		authz:        authz,
	})
}
