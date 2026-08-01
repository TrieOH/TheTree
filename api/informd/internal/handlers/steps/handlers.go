package steps

import (
	"Informd/internal/services"
)

type Handlers struct {
	ops *services.Steps
}

func New(ops *services.Steps) *Handlers { return &Handlers{ops: ops} }

const module = "Informd"
