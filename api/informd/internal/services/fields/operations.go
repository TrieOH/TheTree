package fields

import (
	"Informd/internal/authz"

	"Informd/ports"
	"lib/database"
)

type Operations struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	namespaces ports.NamespaceRepo
	authz      *authz.Service
	// tx opens the transactions multi-step writes run in; threaded from
	// boot, no package-level runner.
	tx database.TxRunner
}

func NewOperations(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	namespaces ports.NamespaceRepo,
	authz *authz.Service,
	tx database.TxRunner,
) *Operations {
	return &Operations{
		forms:      forms,
		steps:      steps,
		fields:     fields,
		namespaces: namespaces,
		authz:      authz,
		tx:         tx,
	}
}
