package projects

import (
	"IdentityX/internal/authz"
	"IdentityX/internal/keys"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	// keys owns the Key-lifecycle: project creation crosses its Ensure
	// seam so every project ships with signing and encryption keys.
	keys     *keys.Manager
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	authz    *authz.Service
}

func NewOperations(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	keysMgr *keys.Manager,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		keys:     keysMgr,
		projects: projects,
		actors:   actors,
		authz:    authz,
	})
}
