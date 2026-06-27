package queries

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	forms      ports.FormsRepo
	steps      ports.StepRepo
	fields     ports.FieldsRepo
	answers    ports.AnswerRepo
	responses  ports.ResponseRepo
	responders ports.ResponderRepo
	namespaces ports.NamespaceRepo
	logger     *zap.Logger
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueries(
	forms ports.FormsRepo,
	steps ports.StepRepo,
	fields ports.FieldsRepo,
	answers ports.AnswerRepo,
	responses ports.ResponseRepo,
	responders ports.ResponderRepo,
	namespaces ports.NamespaceRepo,
	logger *zap.Logger,
	tx database.TxRunner,
	tracer trace.Tracer,
) *Queries {
	return &Queries{
		forms:      forms,
		steps:      steps,
		fields:     fields,
		answers:    answers,
		responses:  responses,
		responders: responders,
		namespaces: namespaces,
		logger:     logger,
		tx:         tx,
		tracer:     tracer,
	}
}
