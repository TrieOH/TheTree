package capabilities

import (
	"IdentityX/internal/services"
)

type Handlers struct {
	ops *services.Capabilities
}

func New(ops *services.Capabilities) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"
