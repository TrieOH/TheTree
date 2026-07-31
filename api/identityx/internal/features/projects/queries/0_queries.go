package queries

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Queries struct {
	projects ports.ProjectRepo
	authz    *authz.Service
}

func NewQueries(
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Queries {
	return errx.MustProvide(&Queries{
		projects: projects,
		authz:    authz,
	})
}
