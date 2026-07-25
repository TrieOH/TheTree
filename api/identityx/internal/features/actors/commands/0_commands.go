package commands

import (
	"IdentityX/ports"
	"lib/errx"
)

type Commands struct {
	actors   ports.ActorRepo
	projects ports.ProjectRepo
}

func NewCommands(
	actors ports.ActorRepo,
	projects ports.ProjectRepo,
) *Commands {
	return errx.MustProvide(&Commands{
		actors:   actors,
		projects: projects,
	})
}
