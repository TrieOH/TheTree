package queries

import (
	"lib/database"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
)

type Queries struct {
	collectors ports.CollectorRepo
	orgs       ports.OrganizationRepo
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewQueries(
	collectors ports.CollectorRepo,
	orgs ports.OrganizationRepo,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		collectors: collectors,
		orgs:       orgs,
		tracer:     tracer,
		tx:         tx,
	}
}
