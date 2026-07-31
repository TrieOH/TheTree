package queries

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Queries struct {
	profiles ports.ProfileRepo
	projects ports.ProjectRepo
	authz    *authz.Service
}

func New(
	profiles ports.ProfileRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Queries {
	return &Queries{
		profiles: profiles,
		projects: projects,
		authz:    authz,
	}
}
