package namespaces

import (
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
}

func NewOperations(
	namespaces ports.NamespaceRepo,
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	answers ports.AnswerRepo,
	responses ports.ResponseRepo,
	responders ports.ResponderRepo,
) *Operations {
	return &Operations{
		namespaces: namespaces,
		forms:      forms,
		steps:      steps,
		fields:     fields,
		answers:    answers,
		responses:  responses,
		responders: responders,
	}
}
