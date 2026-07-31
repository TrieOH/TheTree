package commands

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Commands struct {
	schemas  ports.ProfileSchemaRepo
	projects ports.ProjectRepo
	authz    *authz.Service
}

func New(
	schemas ports.ProfileSchemaRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Commands {
	return &Commands{
		schemas:  schemas,
		projects: projects,
		authz:    authz,
	}
}
