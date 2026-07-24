package queries

import (
	"lib/database"
	"univents/internal/shared/ports"
	ports2 "univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	activities ports.ActivitiesRepository
	editions   ports2.EditionRepo
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewQueries(
	activities ports.ActivitiesRepository,
	editions ports2.EditionRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		activities: activities,
		editions:   editions,
		logger:     logger,
		tracer:     tracer,
		tx:         tx,
	}
}
