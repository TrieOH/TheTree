package queries

import (
	"IdentityX/ports"
)

type Queries struct {
	schemas  ports.ProfileSchemaRepo
	projects ports.ProjectRepo
}

func New(
	schemas ports.ProfileSchemaRepo,
	projects ports.ProjectRepo,
) *Queries {
	return &Queries{
		schemas:  schemas,
		projects: projects,
	}
}
