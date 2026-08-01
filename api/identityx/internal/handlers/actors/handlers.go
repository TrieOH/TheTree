package actors

import (
	"IdentityX/internal/services"
)

type Handlers struct {
	ops *services.Actors
}

func New(ops *services.Actors) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"
