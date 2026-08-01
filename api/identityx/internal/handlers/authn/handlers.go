package authn

import (
	"IdentityX/internal/services"
)

type Handlers struct {
	ops *services.Authn
}

func New(ops *services.Authn) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"
