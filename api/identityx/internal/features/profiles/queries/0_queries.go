package queries

import (
	"IdentityX/ports"
)

type Queries struct {
	profiles ports.ProfileRepo
	projects ports.ProjectRepo
}

func New(
	profiles ports.ProfileRepo,
	projects ports.ProjectRepo,
) *Queries {
	return &Queries{
		profiles: profiles,
		projects: projects,
	}
}
