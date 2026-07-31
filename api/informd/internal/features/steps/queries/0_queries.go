package queries

import (
	"Informd/ports"
)

type Queries struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
}

func NewQueries(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
) *Queries {
	return &Queries{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
	}
}
