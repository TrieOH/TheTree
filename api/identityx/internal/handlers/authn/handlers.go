package authn

import (
	"IdentityX/internal/services"
	"IdentityX/internal/setup"
)

type Handlers struct {
	ops   *services.Authn
	setup *setup.State
}

func New(ops *services.Authn, st *setup.State) *Handlers {
	return &Handlers{ops: ops, setup: st}
}

const module = "IdentityX"
