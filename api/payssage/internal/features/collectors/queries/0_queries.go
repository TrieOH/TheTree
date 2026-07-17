package queries

import (
	"lib/database"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	collectors ports.CollectorRepo
	orgs       ports.OrganizationRepo
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewQueries(
	collectors ports.CollectorRepo,
	orgs ports.OrganizationRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		collectors: collectors,
		orgs:       orgs,
		logger:     logger,
		tracer:     tracer,
		tx:         tx,
	}
}
