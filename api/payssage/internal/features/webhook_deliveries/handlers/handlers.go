package handlers

import (
	"payssage/internal/features/webhook_deliveries/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }
