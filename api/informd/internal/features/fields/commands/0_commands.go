package commands

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type Command struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	namespaces ports.NamespaceRepo
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	namespaces ports.NamespaceRepo,
	tx database.TxRunner,
	tracer trace.Tracer,
) *Command {
	return &Command{
		forms:      forms,
		steps:      steps,
		fields:     fields,
		namespaces: namespaces,
		tx:         tx,
		tracer:     tracer,
	}
}
