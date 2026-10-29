package handlers

import (
	"univents/internal/features/editions/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }
