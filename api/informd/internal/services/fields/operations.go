package fields

import (
	"Informd/ports"
)

type Operations struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	namespaces ports.NamespaceRepo
}

func NewOperations(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	namespaces ports.NamespaceRepo,
) *Operations {
	return &Operations{
		forms:      forms,
		steps:      steps,
		fields:     fields,
		namespaces: namespaces,
	}
}
