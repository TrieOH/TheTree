package commands

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Commands struct {
	profiles ports.ProfileRepo
	schemas  ports.ProfileSchemaRepo
	projects ports.ProjectRepo
	authz    *authz.Service
}

func New(
	profiles ports.ProfileRepo,
	schemas ports.ProfileSchemaRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Commands {
	return &Commands{
		profiles: profiles,
		schemas:  schemas,
		projects: projects,
		authz:    authz,
	}
}
