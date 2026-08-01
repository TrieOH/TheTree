package namespaces

import (
	"Informd/internal/services"
)

type Handlers struct {
	ops *services.Namespaces
}

func New(ops *services.Namespaces) *Handlers { return &Handlers{ops: ops} }

const module = "Informd"
