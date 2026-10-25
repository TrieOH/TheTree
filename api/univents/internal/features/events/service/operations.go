package service

import (
	"lib/database"
	"lib/objectstorage"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/ports"
)

type Deps struct {
	Events ports.EventRepo
	Obj    *objectstorage.Client
	IDX    *idx.Client
	Badges ports.BadgeStaffOps
	Authz  *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
