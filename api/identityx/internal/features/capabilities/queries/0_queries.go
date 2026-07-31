package queries

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Queries struct {
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
	authz        *authz.Service
}

func NewQueries(
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Queries {
	return &Queries{
		capabilities: capabilities,
		projects:     projects,
		authz:        authz,
	}
}
