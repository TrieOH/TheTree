package commands

import (
	"lib/database"
	"lib/errx"
	"payssage/ports"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.opentelemetry.io/otel/trace"
)

type Commands struct {
	river      *river.Client[pgx.Tx]
	events     ports.WebhookEventRepo
	endpoints  ports.WebhookEndpointRepo
	deliveries ports.WebhookDeliveryRepo
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewCommands(
	river *river.Client[pgx.Tx],
	events ports.WebhookEventRepo,
	endpoints ports.WebhookEndpointRepo,
	deliveries ports.WebhookDeliveryRepo,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return errx.MustProvide(&Commands{
		river:      river,
		events:     events,
		endpoints:  endpoints,
		deliveries: deliveries,
		tracer:     tracer,
		tx:         tx,
	})
}
