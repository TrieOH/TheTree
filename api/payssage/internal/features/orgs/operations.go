package orgs

import (
	"payssage/ports"
	idx "sdk/identityx"
)

type Operations struct {
	orgs ports.OrganizationRepo
	idx  *idx.Client
}

func NewOperations(
	orgs ports.OrganizationRepo,
	idx *idx.Client,
) *Operations {
	return &Operations{
		orgs: orgs,
		idx:  idx,
	}
}
