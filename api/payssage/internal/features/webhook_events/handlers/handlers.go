package handlers

import (
	"payssage/internal/features/webhook_events/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
