package commands

import (
	"lib/database"
	"lib/objectstorage"
	"univents/ports"

	idx "sdk/identityx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	events ports.EventRepo
	obj    *objectstorage.Client
	idx    *idx.Client
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewCommands(
	events ports.EventRepo,
	obj *objectstorage.Client,
	idx *idx.Client,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		events: events,
		obj:    obj,
		idx:    idx,
		logger: logger,
		tracer: tracer,
		tx:     tx,
	}
}
