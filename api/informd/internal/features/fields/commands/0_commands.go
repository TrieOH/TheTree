package commands

import (
	"Informd/ports"
)

type Command struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	namespaces ports.NamespaceRepo
}

func NewCommands(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	namespaces ports.NamespaceRepo,
) *Command {
	return &Command{
		forms:      forms,
		steps:      steps,
		fields:     fields,
		namespaces: namespaces,
	}
}
