package webhooks

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.Webhooks
}

func New(ops *services.Webhooks) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
