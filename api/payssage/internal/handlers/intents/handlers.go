package intents

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.Intents
}

func New(ops *services.Intents) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
