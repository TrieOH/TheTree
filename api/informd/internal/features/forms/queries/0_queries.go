package queries

import (
	"Informd/ports"
)

type Queries struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	answers    ports.AnswerRepo
	responses  ports.ResponseRepo
	responders ports.ResponderRepo
	namespaces ports.NamespaceRepo
}

func NewQueries(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	answers ports.AnswerRepo,
	responses ports.ResponseRepo,
	responders ports.ResponderRepo,
	namespaces ports.NamespaceRepo,
) *Queries {
	return &Queries{
		forms:      forms,
		steps:      steps,
		fields:     fields,
		answers:    answers,
		responses:  responses,
		responders: responders,
		namespaces: namespaces,
	}
}
