package checkouts

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Checkouts
}

func New(ops *services.Checkouts) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
