package queries

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
)

type Queries struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	authz    *authz.Service
}

func NewQueries(
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
	authz *authz.Service,
) *Queries {
	return &Queries{
		projects: projects,
		actors:   actors,
		authz:    authz,
	}
}
