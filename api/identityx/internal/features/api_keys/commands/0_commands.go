package commands

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Commands struct {
	hmacSecret   []byte
	actors       ports.ActorRepo
	apiKeys      ports.APIKeysRepo
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
	authz        *authz.Service
}

func NewCommands(
	hmacSecret []byte,
	actors ports.ActorRepo,
	apiKeys ports.APIKeysRepo,
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Commands {
	return &Commands{
		hmacSecret:   hmacSecret,
		actors:       actors,
		apiKeys:      apiKeys,
		capabilities: capabilities,
		projects:     projects,
		authz:        authz,
	}
}
