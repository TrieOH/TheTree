package queries

import (
	"lib/database"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	events ports.EventRepo
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewQueries(
	events ports.EventRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		events: events,
		logger: logger,
		tracer: tracer,
		tx:     tx,
	}
}
