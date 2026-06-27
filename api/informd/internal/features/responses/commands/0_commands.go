package commands

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	responders ports.ResponderRepo
	responses  ports.ResponseRepo
	answers    ports.AnswerRepo
	forms      ports.FormsRepo
	logger     *zap.Logger
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommands(
	responders ports.ResponderRepo,
	responses ports.ResponseRepo,
	answers ports.AnswerRepo,
	forms ports.FormsRepo,
	logger *zap.Logger,
	tx database.TxRunner,
	tracer trace.Tracer,
) *Commands {
	return &Commands{
		responders: responders,
		responses:  responses,
		answers:    answers,
		forms:      forms,
		logger:     logger,
		tx:         tx,
		tracer:     tracer,
	}
}
