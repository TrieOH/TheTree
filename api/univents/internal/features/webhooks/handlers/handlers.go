package handlers

import (
	"univents/internal/features/webhooks/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
