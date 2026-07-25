package queries

import (
	"IdentityX/ports"
)

type Queries struct {
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
}

func NewQueries(
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
) *Queries {
	return &Queries{
		capabilities: capabilities,
		projects:     projects,
	}
}
