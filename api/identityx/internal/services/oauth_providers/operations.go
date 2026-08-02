package oauth_providers

import (
	"IdentityX/internal/authz"
	"IdentityX/ports"
	"lib/errx"
)

type Operations struct {
	providers ports.ProjectOAuthProvidersRepo
	projects  ports.ProjectRepo
	authz     *authz.Service
}

func NewOperations(
	providers ports.ProjectOAuthProvidersRepo,
	projects ports.ProjectRepo,
	authz *authz.Service,
) *Operations {
	return errx.MustProvide(&Operations{
		providers: providers,
		projects:  projects,
		authz:     authz,
	})
}
