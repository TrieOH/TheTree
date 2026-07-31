package commands

import (
	"Informd/ports"
)

type Commands struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
}

func NewCommands(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
) *Commands {
	return &Commands{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
	}
}
