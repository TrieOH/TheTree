package commands

import (
	"lib/database"
	"lib/objectstorage"
	"univents/internal/shared/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	events ports.EventsRepository
	obj    *objectstorage.Client
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewCommands(
	events ports.EventsRepository,
	obj *objectstorage.Client,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		events: events,
		obj:    obj,
		logger: logger,
		tracer: tracer,
		tx:     tx,
	}
}
