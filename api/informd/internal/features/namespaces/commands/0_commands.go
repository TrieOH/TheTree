package commands

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
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
) *CommandService {
	return &CommandService{
		namespaces: projects,
		forms:      forms,
		tx:         tx,
		tracer:     tracer,
	}
}
