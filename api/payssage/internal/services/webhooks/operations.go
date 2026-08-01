package webhooks

import (
	"payssage/ports"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Operations struct {
	river      *river.Client[pgx.Tx]
	events     ports.WebhookEventRepo
	endpoints  ports.WebhookEndpointRepo
	deliveries ports.WebhookDeliveryRepo
}

func NewOperations(
	river *river.Client[pgx.Tx],
	events ports.WebhookEventRepo,
	endpoints ports.WebhookEndpointRepo,
	deliveries ports.WebhookDeliveryRepo,
) *Operations {
	return &Operations{
		river:      river,
		events:     events,
		endpoints:  endpoints,
		deliveries: deliveries,
	}
}
