package queries

import (
	"lib/database"
	ports2 "univents/internal/shared/ports"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	certs    ports.CertificationRepo
	editions ports2.EditionsRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewQueries(
	certs ports.CertificationRepo,
	editions ports2.EditionsRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		certs:    certs,
		editions: editions,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	}
}
