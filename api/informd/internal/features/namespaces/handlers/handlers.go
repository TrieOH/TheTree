package handlers

import (
	"Informd/internal/features/namespaces/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

const module = "Informd"
