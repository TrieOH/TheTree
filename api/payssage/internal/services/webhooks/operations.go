package webhooks

import (
	"payssage/internal/authz"

	"payssage/internal/providers"
	"payssage/ports"

	"lib/database"
	libjobs "lib/jobs"
)

type Operations struct {
	jobs       libjobs.Enqueuer
	events     ports.WebhookEventRepo
	endpoints  ports.WebhookEndpointRepo
	deliveries ports.WebhookDeliveryRepo
	authz      *authz.Service
	registry   providers.PaymentProviders
	// tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	tx database.TxRunner
}

func NewOperations(
	jobs libjobs.Enqueuer,
	events ports.WebhookEventRepo,
	endpoints ports.WebhookEndpointRepo,
	deliveries ports.WebhookDeliveryRepo,
	authz *authz.Service,
	tx database.TxRunner,
	registry providers.PaymentProviders,
) *Operations {
	return &Operations{
		jobs:       jobs,
		events:     events,
		endpoints:  endpoints,
		deliveries: deliveries,
		authz:      authz,
		tx:         tx,
		registry:   registry,
	}
}
