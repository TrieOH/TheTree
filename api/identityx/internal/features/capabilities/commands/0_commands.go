package commands

import (
	"IdentityX/ports"
)

type Commands struct {
	actors       ports.ActorRepo
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
}

func NewCommands(
	actors ports.ActorRepo,
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
) *Commands {
	return &Commands{
		actors:       actors,
		capabilities: capabilities,
		projects:     projects,
	}
}
