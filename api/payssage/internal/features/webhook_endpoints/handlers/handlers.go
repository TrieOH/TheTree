package handlers

import (
	"payssage/internal/features/webhook_endpoints/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
