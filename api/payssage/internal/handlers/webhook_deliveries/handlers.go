package webhook_deliveries

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.WebhookDeliveries
}

func New(ops *services.WebhookDeliveries) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
