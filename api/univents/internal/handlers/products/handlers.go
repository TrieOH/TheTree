package products

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.Products
}

func New(ops *services.Products) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
