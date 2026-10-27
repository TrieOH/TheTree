package service

import (
	"lib/database"
	"payssage/internal/authz"

	"payssage/ports"
	idx "sdk/identityx"
)

type Deps struct {
	Orgs  ports.OrganizationRepo
	IDX   *idx.Client
	Authz *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
