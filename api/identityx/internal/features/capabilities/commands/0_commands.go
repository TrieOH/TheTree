package commands

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Commands struct {
	actors       ports.ActorRepo
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
	authz        *authz.Service
}

func NewCommands(
	actors ports.ActorRepo,
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Commands {
	return &Commands{
		actors:       actors,
		capabilities: capabilities,
		projects:     projects,
		authz:        authz,
	}
}
