package webhook_endpoints

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.WebhookEndpoints
}

func New(ops *services.WebhookEndpoints) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
