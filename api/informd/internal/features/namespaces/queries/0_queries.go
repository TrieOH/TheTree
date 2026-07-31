package queries

import (
	"Informd/ports"
)

type Queries struct {
	namespaces ports.NamespaceRepo
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	answers    ports.AnswerRepo
	responses  ports.ResponseRepo
	responders ports.ResponderRepo
}

func NewQueries(
	namespaces ports.NamespaceRepo,
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	answers ports.AnswerRepo,
	responses ports.ResponseRepo,
	responders ports.ResponderRepo,
) *Queries {
	return &Queries{
		namespaces: namespaces,
		forms:      forms,
		steps:      steps,
		fields:     fields,
		answers:    answers,
		responses:  responses,
		responders: responders,
	}
}
