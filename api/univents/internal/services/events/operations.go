package events

import (
	"lib/objectstorage"
	idx "sdk/identityx"
	"univents/ports"
)

type Operations struct {
	events ports.EventRepo
	obj    *objectstorage.Client
	idx    *idx.Client
}

func NewOperations(
	events ports.EventRepo,
	obj *objectstorage.Client,
	idx *idx.Client,
) *Operations {
	return &Operations{
		events: events,
		obj:    obj,
		idx:    idx,
	}
}
