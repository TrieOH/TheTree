package handlers

import (
	"IdentityX/internal/features/authn/service"
	"IdentityX/internal/setup"
)

type Handlers struct {
	ops   *service.Operations
	setup *setup.State
}

func New(ops *service.Operations, st *setup.State) *Handlers {
	return &Handlers{ops: ops, setup: st}
}
