package queries

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Queries struct {
	schemas  ports.ProfileSchemaRepo
	projects ports.ProjectRepo
	authz    *authz.Service
}

func New(
	schemas ports.ProfileSchemaRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Queries {
	return &Queries{
		schemas:  schemas,
		projects: projects,
		authz:    authz,
	}
}
