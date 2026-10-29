package handlers

import (
	"IdentityX/internal/features/capabilities/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }
