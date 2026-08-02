package orgs

import (
	"payssage/internal/authz"

	"payssage/ports"
	idx "sdk/identityx"
)

type Operations struct {
	orgs  ports.OrganizationRepo
	idx   *idx.Client
	authz *authz.Service
}

func NewOperations(
	orgs ports.OrganizationRepo,
	idx *idx.Client,
	authz *authz.Service,
) *Operations {
	return &Operations{
		orgs:  orgs,
		idx:   idx,
		authz: authz,
	}
}
