package webhooks

import (
	"payssage/internal/authz"

	"payssage/ports"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Operations struct {
	river      *river.Client[pgx.Tx]
	events     ports.WebhookEventRepo
	endpoints  ports.WebhookEndpointRepo
	deliveries ports.WebhookDeliveryRepo
	authz      *authz.Service
}

func NewOperations(
	river *river.Client[pgx.Tx],
	events ports.WebhookEventRepo,
	endpoints ports.WebhookEndpointRepo,
	deliveries ports.WebhookDeliveryRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		river:      river,
		events:     events,
		endpoints:  endpoints,
		deliveries: deliveries,
		authz:      authz,
	}
}
