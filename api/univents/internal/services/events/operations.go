package events

import (
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
}

func NewOperations(
	events ports.EventRepo,
	obj *objectstorage.Client,
	idx *idx.Client,
	authz *authz.Service,
	badges ports.BadgeStaffOps,
) *Operations {
	return &Operations{
		events: events,
		obj:    obj,
		idx:    idx,
		badges: badges,
		authz:  authz,
	}
}
