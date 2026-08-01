package programs

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Programs
}

func New(ops *services.Programs) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
