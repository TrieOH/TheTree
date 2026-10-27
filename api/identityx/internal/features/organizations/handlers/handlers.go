package handlers

import (
	"IdentityX/internal/features/organizations/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"
