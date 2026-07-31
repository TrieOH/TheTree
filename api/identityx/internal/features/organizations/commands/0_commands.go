package commands

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Commands struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	orgs     ports.OrganizationRepo
	authz    *authz.Service
}

func NewCommands(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	orgs ports.OrganizationRepo,
	authz *authz.Service,
) *Commands {
	return errx.MustProvide(&Commands{
		projects: projects,
		actors:   actors,
		orgs:     orgs,
		authz:    authz,
	})
}
