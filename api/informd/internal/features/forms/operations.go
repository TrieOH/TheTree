package forms

import (
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
}

func NewOperations(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	namespaces ports.NamespaceRepo,
	fields ports.FieldsRepo,
	answers ports.AnswerRepo,
	responses ports.ResponseRepo,
	responders ports.ResponderRepo,
) *Operations {
	return &Operations{
		forms:      forms,
		steps:      steps,
		namespaces: namespaces,
		fields:     fields,
		answers:    answers,
		responses:  responses,
		responders: responders,
	}
}
