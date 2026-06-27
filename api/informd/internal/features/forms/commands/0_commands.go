package commands

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
	logger     *zap.Logger
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
	logger *zap.Logger,
	tx database.TxRunner,
	tracer trace.Tracer,
) *Commands {
	return &Commands{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
		logger:     logger,
		tx:         tx,
		tracer:     tracer,
	}
}
