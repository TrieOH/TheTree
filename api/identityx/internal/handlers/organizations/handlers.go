package organizations

import (
	"IdentityX/internal/services"
)

type Handlers struct {
	ops *services.Organizations
}

func New(ops *services.Organizations) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"
