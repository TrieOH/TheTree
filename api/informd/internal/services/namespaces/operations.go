package namespaces

import (
	"Informd/internal/authz"

	"Informd/ports"
	"lib/database"
)

type Operations struct {
	namespaces ports.NamespaceRepo
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	answers    ports.AnswerRepo
	responses  ports.ResponseRepo
	responders ports.ResponderRepo
	authz      *authz.Service
	// tx opens the transactions multi-step writes run in; threaded from
	// boot, no package-level runner.
	tx database.TxRunner
}

func NewOperations(
	namespaces ports.NamespaceRepo,
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	answers ports.AnswerRepo,
	responses ports.ResponseRepo,
	responders ports.ResponderRepo,
	authz *authz.Service,
	tx database.TxRunner,
) *Operations {
	return &Operations{
		namespaces: namespaces,
		forms:      forms,
		steps:      steps,
		fields:     fields,
		answers:    answers,
		responses:  responses,
		responders: responders,
		authz:      authz,
		tx:         tx,
	}
}
