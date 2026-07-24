package commands

import (
	"lib/database"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	certs    ports.CertificationRepo
	editions ports.EditionRepo
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewCommands(
	certs ports.CertificationRepo,
	editions ports.EditionRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		certs:    certs,
		editions: editions,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	}
}
