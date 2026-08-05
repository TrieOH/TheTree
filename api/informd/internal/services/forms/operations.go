package forms

import (
	"Informd/internal/authz"

	"Informd/ports"
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
	}
}
