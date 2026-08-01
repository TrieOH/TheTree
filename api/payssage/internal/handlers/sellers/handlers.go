package sellers

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.Sellers
}

func New(ops *services.Sellers) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
