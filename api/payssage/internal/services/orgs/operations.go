package orgs

import (
	"lib/database"
	"payssage/internal/authz"

	"payssage/ports"
	idx "sdk/identityx"
)

type Operations struct {
	orgs  ports.OrganizationRepo
	idx   *idx.Client
	authz *authz.Service
	// tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	tx database.TxRunner
}

func NewOperations(
	orgs ports.OrganizationRepo,
	idx *idx.Client,
	authz *authz.Service,
	tx database.TxRunner,
) *Operations {
	return &Operations{
		orgs:  orgs,
		idx:   idx,
		authz: authz,
		tx:    tx,
	}
}
