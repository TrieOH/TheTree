package queries

import (
	"lib/database"
	"univents/internal/shared/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	events   ports.EventsRepository
	editions ports.EditionsRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewQueries(
	events ports.EventsRepository,
	editions ports.EditionsRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		events:   events,
		editions: editions,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	}
}
