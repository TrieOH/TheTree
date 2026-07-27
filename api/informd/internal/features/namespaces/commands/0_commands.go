package commands

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type Commands struct {
	namespaces ports.NamespaceRepo
	forms      ports.FormsRepo
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	projects ports.NamespaceRepo,
	forms ports.FormsRepo,
	tx database.TxRunner,
	tracer trace.Tracer,
) *Commands {
	return &Commands{
		namespaces: projects,
		forms:      forms,
		tx:         tx,
		tracer:     tracer,
	}
}
