package handlers

import (
	"Informd/internal/features/steps/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

const module = "Informd"
