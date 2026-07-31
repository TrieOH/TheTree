package commands

import (
	"lib/errx"
	"payssage/ports"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Commands struct {
	river      *river.Client[pgx.Tx]
	events     ports.WebhookEventRepo
	endpoints  ports.WebhookEndpointRepo
	deliveries ports.WebhookDeliveryRepo
}

func NewCommands(
	river *river.Client[pgx.Tx],
	events ports.WebhookEventRepo,
	endpoints ports.WebhookEndpointRepo,
	deliveries ports.WebhookDeliveryRepo,
) *Commands {
	return errx.MustProvide(&Commands{
		river:      river,
		events:     events,
		endpoints:  endpoints,
		deliveries: deliveries,
	})
}
