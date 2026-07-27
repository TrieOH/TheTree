package commands

import (
	"lib/objectstorage"
	"univents/ports"

	idx "sdk/identityx"
)

type Commands struct {
	events ports.EventRepo
	obj    *objectstorage.Client
	idx    *idx.Client
}

func NewCommands(
	events ports.EventRepo,
	obj *objectstorage.Client,
	idx *idx.Client,
) *Commands {
	return &Commands{
		events: events,
		obj:    obj,
		idx:    idx,
	}
}
