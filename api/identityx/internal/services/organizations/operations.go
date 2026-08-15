package organizations

import (
	"IdentityX/internal/authz"
	"IdentityX/internal/keys"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	orgs     ports.OrganizationRepo
	// keys owns the Key-lifecycle: org-created projects cross its Ensure
	// seam so they ship with keys instead of being token-broken until the
	// next boot.
	keys  *keys.Manager
	authz *authz.Service
}

func NewOperations(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	orgs ports.OrganizationRepo,
	keysMgr *keys.Manager,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		projects: projects,
		actors:   actors,
		orgs:     orgs,
		keys:     keysMgr,
		authz:    authz,
	})
}
