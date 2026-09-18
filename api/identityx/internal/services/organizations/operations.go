package organizations

import (
	"IdentityX/internal/authz"
	"IdentityX/internal/keys"
	"IdentityX/ports"
	"lib/database"
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
	// tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	tx database.TxRunner
}

func NewOperations(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	orgs ports.OrganizationRepo,
	keysMgr *keys.Manager,
	authz *authz.Service,
	tx database.TxRunner,
) *Operations {
	return errx.MustProvide(&Operations{
		projects: projects,
		actors:   actors,
		orgs:     orgs,
		keys:     keysMgr,
		authz:    authz,
		tx:       tx,
	})
}
