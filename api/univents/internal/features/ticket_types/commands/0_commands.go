package commands

import (
	"lib/database"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	events      ports.EventRepo
	editions    ports.EditionRepo
	ticketTypes ports.TicketTypeRepo
	logger      *zap.Logger
	tracer      trace.Tracer
	tx          database.TxRunner
}

func NewCommands(
	events ports.EventRepo,
	editions ports.EditionRepo,
	ticketTypes ports.TicketTypeRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		events:      events,
		editions:    editions,
		ticketTypes: ticketTypes,
		logger:      logger,
		tracer:      tracer,
		tx:          tx,
	}
}
