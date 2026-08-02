package steps

import (
	"Informd/internal/authz"

	"Informd/ports"
)

type Operations struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
	authz      *authz.Service
}

func NewOperations(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
		authz:      authz,
	}
}
