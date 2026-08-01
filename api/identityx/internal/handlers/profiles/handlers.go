package profiles

import (
	"IdentityX/internal/services"
)

type Handlers struct {
	ops *services.Profiles
}

func New(ops *services.Profiles) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"
