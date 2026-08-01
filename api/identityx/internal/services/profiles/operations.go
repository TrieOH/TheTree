package profiles

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	profiles ports.ProfileRepo
	schemas  ports.ProfileSchemaRepo
	projects ports.ProjectRepo
	authz    *authz.Service
}

func NewOperations(
	profiles ports.ProfileRepo,
	schemas ports.ProfileSchemaRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		profiles: profiles,
		schemas:  schemas,
		projects: projects,
		authz:    authz,
	})
}
