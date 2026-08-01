package organizations

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	orgs     ports.OrganizationRepo
	authz    *authz.Service
}

func NewOperations(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	orgs ports.OrganizationRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		projects: projects,
		actors:   actors,
		orgs:     orgs,
		authz:    authz,
	})
}
