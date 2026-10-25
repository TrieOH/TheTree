package service

import (
	"Informd/internal/authz"

	"Informd/ports"
	"lib/database"
)

type Deps struct {
	Forms      ports.FormsRepo
	Steps      ports.StepRepo
	Fields     ports.FieldsRepo
	Namespaces ports.NamespaceRepo
	Authz      *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded from
	// boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
