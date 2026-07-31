package commands

import (
	"lib/errx"
	"payssage/ports"
	idx "sdk/identityx"
)

type Commands struct {
	orgs ports.OrganizationRepo
	idx  *idx.Client
}

func NewCommands(
	orgs ports.OrganizationRepo,
	idx *idx.Client,
) *Commands {
	return errx.MustProvide(&Commands{
		orgs: orgs,
		idx:  idx,
	})
}
