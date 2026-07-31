package queries

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Queries struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	orgs     ports.OrganizationRepo
	authz    *authz.Service
}

func NewQueries(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	orgs ports.OrganizationRepo,
	authz *authz.Service,
) *Queries {
	return errx.MustProvide(&Queries{
		projects: projects,
		actors:   actors,
		orgs:     orgs,
		authz:    authz,
	})
}
