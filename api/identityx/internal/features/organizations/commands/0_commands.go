package commands

import (
	"IdentityX/ports"
	"lib/errx"
)

type Commands struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	orgs     ports.OrganizationRepo
}

func NewCommands(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	orgs ports.OrganizationRepo,
) *Commands {
	return errx.MustProvide(&Commands{
		projects: projects,
		actors:   actors,
		orgs:     orgs,
	})
}
