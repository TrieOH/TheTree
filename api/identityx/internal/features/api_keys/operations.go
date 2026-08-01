package api_keys

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	hmacSecret   []byte
	actors       ports.ActorRepo
	apiKeys      ports.APIKeysRepo
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
	authz        *authz.Service
}

func NewOperations(
	hmacSecret []byte,
	actors ports.ActorRepo,
	apiKeys ports.APIKeysRepo,
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		hmacSecret:   hmacSecret,
		actors:       actors,
		apiKeys:      apiKeys,
		capabilities: capabilities,
		projects:     projects,
		authz:        authz,
	})
}
