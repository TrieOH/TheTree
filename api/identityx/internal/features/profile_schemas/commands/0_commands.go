package commands

import (
	"IdentityX/ports"
)

type Commands struct {
	schemas  ports.ProfileSchemaRepo
	projects ports.ProjectRepo
}

func New(
	schemas ports.ProfileSchemaRepo,
	projects ports.ProjectRepo,
) *Commands {
	return &Commands{
		schemas:  schemas,
		projects: projects,
	}
}
