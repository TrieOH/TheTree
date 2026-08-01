package editions

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Editions
}

func New(ops *services.Editions) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
