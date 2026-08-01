package oauth

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.OAuth
}

func New(ops *services.OAuth) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
