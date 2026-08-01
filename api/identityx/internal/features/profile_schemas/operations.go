package profile_schemas

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	schemas  ports.ProfileSchemaRepo
	projects ports.ProjectRepo
	authz    *authz.Service
}

func NewOperations(
	schemas ports.ProfileSchemaRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		schemas:  schemas,
		projects: projects,
		authz:    authz,
	})
}
