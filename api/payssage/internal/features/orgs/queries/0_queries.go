package queries

import (
	"lib/errx"
	"payssage/ports"
	idx "sdk/identityx"
)

type Queries struct {
	orgs ports.OrganizationRepo
	idx  *idx.Client
}

func NewQueries(
	orgs ports.OrganizationRepo,
	idx *idx.Client,
) *Queries {
	return errx.MustProvide(&Queries{
		orgs: orgs,
		idx:  idx,
	})
}
