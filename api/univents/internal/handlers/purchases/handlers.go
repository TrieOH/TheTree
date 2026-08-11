package purchases

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Purchases
}

func New(ops *services.Purchases) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
