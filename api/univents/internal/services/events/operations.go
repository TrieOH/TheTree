package events

import (
	"lib/database"
	"lib/objectstorage"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	events ports.EventRepo
	obj    *objectstorage.Client
	idx    *idx.Client
	badges ports.BadgeStaffOps
	authz  *authz.Service
	// tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	tx database.TxRunner
}

func NewOperations(
	events ports.EventRepo,
	obj *objectstorage.Client,
	idx *idx.Client,
	authz *authz.Service,
	badges ports.BadgeStaffOps,
	tx database.TxRunner,
) *Operations {
	return &Operations{
		events: events,
		obj:    obj,
		idx:    idx,
		badges: badges,
		authz:  authz,
		tx:     tx,
	}
}
