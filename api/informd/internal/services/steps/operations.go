package steps

import (
	"Informd/ports"
)

type Operations struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
}

func NewOperations(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
) *Operations {
	return &Operations{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
	}
}
