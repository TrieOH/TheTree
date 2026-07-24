package queries

import (
	"lib/database"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	editions    ports.EditionRepo
	ticketTypes ports.TicketTypeRepo
	logger      *zap.Logger
	tracer      trace.Tracer
	tx          database.TxRunner
}

func NewQueries(
	editions ports.EditionRepo,
	ticketTypes ports.TicketTypeRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		editions:    editions,
		ticketTypes: ticketTypes,
		logger:      logger,
		tracer:      tracer,
		tx:          tx,
	}
}
