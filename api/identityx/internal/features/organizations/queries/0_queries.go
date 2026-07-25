package queries

import (
	"IdentityX/ports"
	"lib/errx"
)

type Queries struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	orgs     ports.OrganizationRepo
}

func NewQueries(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return errx.MustProvide(&Queries{
		projects: projects,
		actors:   actors,
		orgs:     orgs,
	})
}
