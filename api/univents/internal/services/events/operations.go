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
	authz  *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	obj *objectstorage.Client,
	idx *idx.Client,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events: events,
		obj:    obj,
		idx:    idx,
		authz:  authz,
	}
}
