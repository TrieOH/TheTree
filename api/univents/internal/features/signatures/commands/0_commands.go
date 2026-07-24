package commands

import (
	"lib/database"
	"lib/objectstorage"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	signatures ports.SignatureRepo
	editions   ports.EditionRepo
	obj        *objectstorage.Client
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewCommands(
	signatures ports.SignatureRepo,
	editions ports.EditionRepo,
	obj *objectstorage.Client,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		signatures: signatures,
		editions:   editions,
		obj:        obj,
		logger:     logger,
		tracer:     tracer,
		tx:         tx,
	}
}
