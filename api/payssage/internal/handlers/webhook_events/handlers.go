package webhook_events

import (
	"payssage/internal/services"
)

type Handlers struct {
	ops *services.WebhookEvents
}

func New(ops *services.WebhookEvents) *Handlers { return &Handlers{ops: ops} }

const module = "Payssage"
