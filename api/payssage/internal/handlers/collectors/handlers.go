package collectors

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.Collectors
}

func New(ops *services.Collectors) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
