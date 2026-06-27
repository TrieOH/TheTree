package commands

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Command struct {
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
) *Command {
	return &Command{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
		logger:     logger,
		tx:         tx,
		tracer:     tracer,
	}
}
