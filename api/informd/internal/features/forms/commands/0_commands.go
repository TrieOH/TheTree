package commands

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
		tx:         tx,
		tracer:     tracer,
	}
}
