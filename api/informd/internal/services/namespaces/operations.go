package namespaces

import (
	"Informd/internal/authz"

	"Informd/ports"
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
	}
}
