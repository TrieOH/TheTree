package profiles

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	profiles ports.ProfileRepo
	schemas  ports.ProfileSchemaRepo
	actors   ports.ActorRepo
	authz    *authz.Service
}

func NewOperations(
	profiles ports.ProfileRepo,
	schemas ports.ProfileSchemaRepo,
	actors ports.ActorRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		profiles: profiles,
		schemas:  schemas,
		actors:   actors,
		authz:    authz,
	})
}
