package commands

import (
	"IdentityX/ports"
)

type Commands struct {
	profiles ports.ProfileRepo
	schemas  ports.ProfileSchemaRepo
	projects ports.ProjectRepo
}

func New(
	profiles ports.ProfileRepo,
	schemas ports.ProfileSchemaRepo,
	projects ports.ProjectRepo,
) *Commands {
	return &Commands{
		profiles: profiles,
		schemas:  schemas,
		projects: projects,
	}
}
