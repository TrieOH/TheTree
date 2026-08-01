package profile_schemas

import (
	"IdentityX/internal/services"
)

type Handlers struct {
	ops *services.ProfileSchemas
}

func New(ops *services.ProfileSchemas) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"
