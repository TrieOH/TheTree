package commands

import (
	"IdentityX/ports"
)

type Commands struct {
	hmacSecret   []byte
	actors       ports.ActorRepo
	apiKeys      ports.APIKeysRepo
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
}

func NewCommands(
	hmacSecret []byte,
	actors ports.ActorRepo,
	apiKeys ports.APIKeysRepo,
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
) *Commands {
	return &Commands{
		hmacSecret:   hmacSecret,
		actors:       actors,
		apiKeys:      apiKeys,
		capabilities: capabilities,
		projects:     projects,
	}
}
