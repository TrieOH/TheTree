package webhooks

import (
	"payssage/internal/authz"

	"payssage/ports"

	"lib/database"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Operations struct {
	river      *river.Client[pgx.Tx]
	events     ports.WebhookEventRepo
	endpoints  ports.WebhookEndpointRepo
	deliveries ports.WebhookDeliveryRepo
	authz      *authz.Service
	// tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	tx database.TxRunner
}

func NewOperations(
	river *river.Client[pgx.Tx],
	events ports.WebhookEventRepo,
	endpoints ports.WebhookEndpointRepo,
	deliveries ports.WebhookDeliveryRepo,
	authz *authz.Service,
	tx database.TxRunner,
) *Operations {
	return &Operations{
		river:      river,
		events:     events,
		endpoints:  endpoints,
		deliveries: deliveries,
		authz:      authz,
		tx:         tx,
	}
}
