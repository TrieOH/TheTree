package commands

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	namespaces ports.NamespaceRepo
	forms      ports.FormsRepo
	logger     *zap.Logger
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	projects ports.NamespaceRepo,
	forms ports.FormsRepo,
	logger *zap.Logger,
	tx database.TxRunner,
	tracer trace.Tracer,
) *Commands {
	return &Commands{
		namespaces: projects,
		forms:      forms,
		logger:     logger,
		tx:         tx,
		tracer:     tracer,
	}
}
