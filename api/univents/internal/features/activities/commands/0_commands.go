package commands

import (
	"lib/database"
	"univents/internal/shared/ports"
	ports2 "univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	activities ports.ActivitiesRepository
	editions   ports2.EditionRepo
	certs      ports2.CertificationRepo
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewCommands(
	activities ports.ActivitiesRepository,
	editions ports2.EditionRepo,
	certs ports2.CertificationRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		activities: activities,
		editions:   editions,
		certs:      certs,
		logger:     logger,
		tracer:     tracer,
		tx:         tx,
	}
}
