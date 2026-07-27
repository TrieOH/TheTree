package commands

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type Commands struct {
	responders ports.ResponderRepo
	responses  ports.ResponseRepo
	answers    ports.AnswerRepo
	forms      ports.FormsRepo
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	responders ports.ResponderRepo,
	responses ports.ResponseRepo,
	answers ports.AnswerRepo,
	forms ports.FormsRepo,
	tx database.TxRunner,
	tracer trace.Tracer,
) *Commands {
	return &Commands{
		responders: responders,
		responses:  responses,
		answers:    answers,
		forms:      forms,
		tx:         tx,
		tracer:     tracer,
	}
}
