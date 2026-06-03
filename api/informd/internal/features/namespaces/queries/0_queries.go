package queries

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type Queries struct {
	namespaces ports.NamespaceRepo
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	answers    ports.AnswerRepo
	responses  ports.ResponseRepo
	responders ports.ResponderRepo
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueries(
	namespaces ports.NamespaceRepo,
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	answers ports.AnswerRepo,
	responses ports.ResponseRepo,
	responders ports.ResponderRepo,
	tx database.TxRunner,
	tracer trace.Tracer,
) *Queries {
	return &Queries{
		namespaces: namespaces,
		forms:      forms,
		steps:      steps,
		fields:     fields,
		answers:    answers,
		responses:  responses,
		responders: responders,
		tx:         tx,
		tracer:     tracer,
	}
}
