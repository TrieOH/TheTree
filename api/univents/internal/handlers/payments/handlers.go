package payments

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Payments
}

func New(ops *services.Payments) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
