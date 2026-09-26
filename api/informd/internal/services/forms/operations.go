package forms

import (
	"Informd/internal/authz"

	"Informd/ports"
	"lib/database"
)

type Operations struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	namespaces ports.NamespaceRepo
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
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
	fields ports.FieldsRepo,
	answers ports.AnswerRepo,
	responses ports.ResponseRepo,
	responders ports.ResponderRepo,
	authz *authz.Service,
	tx database.TxRunner,
) *Operations {
	return &Operations{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
		fields:     fields,
		answers:    answers,
		responses:  responses,
		responders: responders,
		authz:      authz,
		tx:         tx,
	}
}
