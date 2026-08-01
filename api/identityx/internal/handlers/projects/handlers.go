package projects

import (
	"IdentityX/internal/services"
)

type Handlers struct {
	ops *services.Projects
}

func New(ops *services.Projects) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"
